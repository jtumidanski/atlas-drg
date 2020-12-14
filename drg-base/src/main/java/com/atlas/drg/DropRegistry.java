package com.atlas.drg;

import java.util.ArrayList;
import java.util.List;
import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.atomic.AtomicInteger;

import com.atlas.drg.model.Drop;

public class DropRegistry {
   private static final Object lock = new Object();

   private static final Object registryLock = new Object();

   private static volatile DropRegistry instance;

   private final Map<Integer, Drop> dropMap;

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
   }

   public Drop createDrop(int itemId, int quantity, int type, int x, int y, int ownerId, Integer ownerPartyId, long dropTime,
                          int dropperId, int dropperX, int dropperY, boolean playerDrop, byte mod) {
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
         result = new Drop(currentUniqueId, itemId, quantity, type, x, y, ownerId, ownerPartyId, dropTime, dropperId, dropperX,
               dropperY, playerDrop, mod);
         dropMap.put(currentUniqueId, result);
      }
      return result;
   }
}
