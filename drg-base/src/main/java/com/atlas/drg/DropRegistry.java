package com.atlas.drg;

import java.util.ArrayList;
import java.util.Collection;
import java.util.Collections;
import java.util.HashSet;
import java.util.List;
import java.util.Map;
import java.util.Optional;
import java.util.Set;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.atomic.AtomicInteger;
import java.util.stream.Collectors;

import com.atlas.drg.model.Drop;
import com.atlas.drg.model.DropStatus;
import com.atlas.drg.model.MapKey;

public class DropRegistry {
   private static final Object lock = new Object();

   private static final Object registryLock = new Object();

   private static volatile DropRegistry instance;

   private final Map<Integer, Drop> dropMap;

   private final Map<MapKey, Set<Integer>> dropsInMapMap;

   private final Map<Integer, Integer> dropReservation;

   private final AtomicInteger runningUniqueId = new AtomicInteger(1000000001);

   public static DropRegistry getInstance() {
      DropRegistry result = instance;
      if (result == null) {
         synchronized (lock) {
            result = instance;
            if (result == null) {
               result = new DropRegistry();
               instance = result;
            }
         }
      }
      return result;
   }

   private DropRegistry() {
      dropMap = new ConcurrentHashMap<>();
      dropsInMapMap = new ConcurrentHashMap<>();
      dropReservation = new ConcurrentHashMap<>();
   }

   public Drop createDrop(int worldId, int channelId, int mapId, int itemId, int quantity, int meso, int type, int x, int y,
                          int ownerId, Integer ownerPartyId, long dropTime, int dropperId, int dropperX, int dropperY,
                          boolean playerDrop, byte mod) {
      Integer currentUniqueId;
      synchronized (registryLock) {
         List<Integer> existingIds = new ArrayList<>(dropMap.keySet());
         do {
            if ((currentUniqueId = runningUniqueId.incrementAndGet()) >= 2000000000) {
               runningUniqueId.set(currentUniqueId = 1000000001);
            }
         } while (existingIds.contains(currentUniqueId));
      }
      Drop result;
      synchronized (currentUniqueId) {
         result = new Drop(currentUniqueId, worldId, channelId, mapId, itemId, quantity, meso, type, x, y, ownerId, ownerPartyId,
               dropTime, dropperId, dropperX, dropperY, playerDrop, mod, DropStatus.AVAILABLE);
         dropMap.put(currentUniqueId, result);
         MapKey mapKey = new MapKey(worldId, channelId, mapId);
         synchronized (mapKey) {
            if (!dropsInMapMap.containsKey(mapKey)) {
               dropsInMapMap.put(mapKey, new HashSet<>());
            }
            dropsInMapMap.get(mapKey).add(currentUniqueId);
         }
      }
      return result;
   }

   public List<Drop> getDropsForMap(int worldId, int channelId, int mapId) {
      MapKey mapKey = new MapKey(worldId, channelId, mapId);
      synchronized (mapKey) {
         return dropsInMapMap.get(mapKey).stream()
               .map(dropMap::get)
               .collect(Collectors.toUnmodifiableList());
      }
   }

   public Collection<Drop> getDrops() {
      return Collections.unmodifiableCollection(dropMap.values());
   }

   public Optional<Drop> removeDrop(Integer uniqueId) {
      synchronized (uniqueId) {
         Drop drop = null;
         if (dropMap.containsKey(uniqueId)) {
            drop = dropMap.get(uniqueId);
            dropMap.remove(uniqueId);

            MapKey mapKey = new MapKey(drop.worldId(), drop.channelId(), drop.mapId());
            synchronized (mapKey) {
               if (dropsInMapMap.containsKey(mapKey)) {
                  dropsInMapMap.get(mapKey).remove(uniqueId);
               }
            }
         }
         dropReservation.remove(uniqueId);
         return Optional.ofNullable(drop);
      }
   }

   public Optional<Drop> getDrop(Integer uniqueId) {
      synchronized (uniqueId) {
         if (dropMap.containsKey(uniqueId)) {
            return Optional.of(dropMap.get(uniqueId));
         }
         return Optional.empty();
      }
   }

   public Optional<Drop> reserveDrop(Integer uniqueId, int characterId) {
      synchronized (uniqueId) {
         if (dropMap.containsKey(uniqueId)) {
            Drop drop = dropMap.get(uniqueId);
            if (drop.status().equals(DropStatus.AVAILABLE)) {
               Drop updatedDrop = drop.reserve();
               dropMap.put(uniqueId, updatedDrop);
               dropReservation.put(uniqueId, characterId);
               return Optional.of(updatedDrop);
            } else {
               if (dropReservation.containsKey(uniqueId) && dropReservation.get(uniqueId).equals(characterId)) {
                  return Optional.of(drop);
               }
            }
         }
      }
      return Optional.empty();
   }
}
