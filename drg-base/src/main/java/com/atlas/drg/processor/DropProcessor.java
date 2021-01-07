package com.atlas.drg.processor;

import java.awt.*;
import java.util.List;
import java.util.Optional;
import java.util.Random;
import java.util.concurrent.CompletableFuture;
import java.util.stream.Collectors;
import java.util.stream.IntStream;

import com.app.rest.util.RestResponseUtil;
import com.atlas.dis.rest.attribute.MonsterDropAttributes;
import com.atlas.drg.DropRegistry;
import com.atlas.drg.event.producer.DropEventProducer;
import com.atlas.drg.event.producer.DropExpiredEventProducer;
import com.atlas.drg.event.producer.DropPickedUpEventProducer;
import com.atlas.drg.event.producer.DropReservationEventProducer;
import com.atlas.drg.model.Drop;
import com.atlas.drg.model.MonsterDrop;
import com.atlas.mis.attribute.DropPositionInputAttributes;
import com.atlas.mis.attribute.MapPointAttributes;
import com.atlas.mis.builder.DropPositionInputAttributesBuilder;
import com.atlas.mis.constant.RestConstants;
import com.atlas.shared.rest.UriBuilder;

import builder.ResultObjectBuilder;
import rest.DataBody;
import rest.DataContainer;

public final class DropProcessor {
   private DropProcessor() {
   }

   public static void createDrops(int worldId, int channelId, int mapId, int monsterUniqueId, int monsterId, int x, int y,
                                  int killerId) {
      // TODO determine type of drop
      //    monster is explosive? 3
      //    monster has ffa loot? 2
      //    killer is in party? 1
      byte dropType = 0;

      getMonsterDropStream(monsterId)
            .thenApply(drops -> getSuccessfulDrops(killerId, drops))
            .thenAccept(drops -> IntStream.range(0, drops.size())
                  .forEach(i -> createDrop(worldId, channelId, mapId, i + 1, monsterUniqueId, x, y, killerId, dropType,
                        drops.get(i))));
   }

   protected static CompletableFuture<List<MonsterDrop>> getMonsterDropStream(int monsterId) {
      return UriBuilder.service(com.atlas.drg.constant.RestConstants.SERVICE)
            .path("monsters")
            .path("drops")
            .queryParam("monsterId", monsterId)
            .getAsyncRestClient(MonsterDropAttributes.class)
            .get()
            .thenApply(RestResponseUtil::result)
            .thenApply(DataContainer::dataList)
            .thenApply(DropProcessor::getMonsterDropsFromBody);
   }

   protected static List<MonsterDrop> getSuccessfulDrops(int killerId, List<MonsterDrop> allDrops) {
      return allDrops.stream()
            .filter(monsterDrop -> evaluateSuccess(killerId, monsterDrop))
            .collect(Collectors.toList());
   }

   protected static List<MonsterDrop> getMonsterDropsFromBody(List<DataBody<MonsterDropAttributes>> body) {
      return body.stream()
            .map(ModelFactory::createMonsterDrop)
            .collect(Collectors.toList());
   }

   /**
    * Evaluates the success or failure of the drop.
    *
    * @param killerId    the character whose rates will be considered
    * @param monsterDrop the drop being evaluated
    * @return true if the drop should drop
    */
   protected static boolean evaluateSuccess(int killerId, MonsterDrop monsterDrop) {
      //TODO evaluate card rate for killer, whether it's meso or drop.
      int chance = (int) Math.min((float) monsterDrop.chance() * 1, Integer.MAX_VALUE);
      return new Random().nextInt(999999) < chance;
   }

   protected static void createDrop(int worldId, int channelId, int mapId, int index, int monsterUniqueId, int x, int y,
                                    int killerId,
                                    byte dropType, MonsterDrop monsterDrop) {
      int factor;
      if (dropType == 3) {
         factor = 40;
      } else {
         factor = 25;
      }
      int newX = x + ((index % 2 == 0) ? (factor * ((index + 1) / 2)) : -(factor * (index / 2)));
      Point position = new Point(newX, y);
      if (monsterDrop.itemId() == 0) {
         spawnMeso(worldId, channelId, mapId, monsterUniqueId, x, y, killerId, dropType, monsterDrop, position);
      } else {
         spawnItem(worldId, channelId, mapId, monsterDrop.itemId(), monsterUniqueId, x, y, killerId, dropType, monsterDrop,
               position);
      }
   }

   protected static void spawnMeso(int worldId, int channelId, int mapId, int monsterUniqueId, int x, int y, int killerId,
                                   byte dropType,
                                   MonsterDrop drop, Point position) {
      int mesos = new Random().nextInt(drop.maximumQuantity() - drop.minimumQuantity()) + drop.minimumQuantity();
      if (mesos > 0) {
         //TODO apply characters meso buff.
         spawnDrop(worldId, channelId, mapId, 0, 0, mesos, position.x, position.y, x, y, monsterUniqueId, killerId,
               false, dropType);
      }
   }

   protected static void spawnItem(int worldId, int channelId, int mapId, int itemId, int monsterUniqueId, int x, int y,
                                   int killerId,
                                   byte dropType,
                                   MonsterDrop drop, Point position) {
      int quantity;
      if (drop.maximumQuantity() == 1) {
         quantity = 1;
      } else {
         quantity = new Random().nextInt(drop.maximumQuantity() - drop.minimumQuantity()) + drop.minimumQuantity();
      }
      spawnDrop(worldId, channelId, mapId, itemId, quantity, 0, position.x, position.y, x, y, monsterUniqueId, killerId,
            false, dropType);
   }

   protected static void spawnDrop(int worldId, int channelId, int mapId, int itemId, int quantity, int meso, int itemX, int itemY,
                                   int monsterX, int monsterY, int monsterUniqueId, int killerId, boolean playerDrop,
                                   byte dropType) {
      calculateDropPosition(mapId, itemX, itemY, monsterX, monsterY)
            .thenCompose(position -> calculateDropPosition(mapId, position.x, position.y, position.x, position.y))
            .thenApply(position -> DropRegistry.getInstance().createDrop(worldId, channelId, mapId, itemId, quantity, meso,
                  dropType, position.x, position.y, killerId, null, System.currentTimeMillis(), monsterUniqueId, monsterX,
                  monsterY, playerDrop, (byte) 1))
            .thenAccept(drop -> DropEventProducer.createDrop(worldId, channelId, mapId, drop));
   }

   protected static CompletableFuture<Point> calculateDropPosition(int mapId, int initialX, int initialY, int fallbackX,
                                                                   int fallbackY) {
      return UriBuilder.service(RestConstants.SERVICE)
            .pathParam("maps", mapId)
            .path("dropPosition")
            .getAsyncRestClient(MapPointAttributes.class)
            .create(new ResultObjectBuilder(DropPositionInputAttributes.class, 0)
                  .setAttribute(new DropPositionInputAttributesBuilder()
                        .setInitialX(initialX)
                        .setInitialY(initialY)
                        .setFallbackX(fallbackX)
                        .setFallbackY(fallbackY)
                  )
                  .inputObject()
            )
            .thenApply(RestResponseUtil::result)
            .thenApply(DataContainer::data)
            .thenApply(Optional::get)
            .thenApply(body -> new Point(body.getAttributes().x(), body.getAttributes().y()))
            .exceptionally(fh -> new Point(fallbackX, fallbackY));
   }

   public static void destroyAll() {
      DropRegistry.getInstance().getDrops().forEach(DropProcessor::destroyDrop);
   }

   public static void destroyDrop(Drop drop) {
      DropRegistry.getInstance().removeDrop(drop.id());
      DropExpiredEventProducer.expireDrop(drop.worldId(), drop.channelId(), drop.mapId(), drop.id());
   }

   public static void reserveDrop(int dropId, int characterId) {
      DropRegistry.getInstance().reserveDrop(dropId, characterId)
            .ifPresentOrElse(
                  drop -> DropReservationEventProducer.reservationSuccess(dropId, characterId),
                  () -> DropReservationEventProducer.reservationFailure(dropId, characterId)
            );
   }

   public static void pickupDrop(int dropId, int characterId) {
      DropRegistry.getInstance()
            .removeDrop(dropId)
            .ifPresent(drop -> DropPickedUpEventProducer.emit(dropId, characterId, drop.mapId()));
   }

   public static void cancelDropReservation(int dropId, int characterId) {
      DropRegistry.getInstance().cancelDropReservation(dropId, characterId);
      DropReservationEventProducer.reservationFailure(dropId, characterId);
   }
}
