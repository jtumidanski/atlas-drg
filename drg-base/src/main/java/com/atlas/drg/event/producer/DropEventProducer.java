package com.atlas.drg.event.producer;

import com.atlas.drg.EventProducerRegistry;
import com.atlas.drg.constant.EventConstants;
import com.atlas.drg.event.DropEvent;
import com.atlas.drg.model.Drop;

public final class DropEventProducer {
   private DropEventProducer() {
   }

   public static void createDrop(int worldId, int channelId, int mapId, Drop drop) {
      EventProducerRegistry.getInstance().send(EventConstants.TOPIC_DROP_EVENT, mapId,
            new DropEvent(worldId, channelId, mapId, drop.id(), drop.itemId(), drop.quantity(), drop.meso(), drop.type(), drop.x(),
                  drop.y(), drop.ownerId(), drop.ownerPartyId(), drop.dropTime(), drop.dropperId(), drop.dropperX(),
                  drop.dropperY(), drop.playerDrop(), drop.mod()));
   }
}
