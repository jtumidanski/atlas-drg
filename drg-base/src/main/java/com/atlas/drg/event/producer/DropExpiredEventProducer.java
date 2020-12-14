package com.atlas.drg.event.producer;

import com.atlas.drg.EventProducerRegistry;
import com.atlas.drg.constant.EventConstants;
import com.atlas.drg.event.DropExpiredEvent;

public final class DropExpiredEventProducer {
   private DropExpiredEventProducer() {
   }

   public static void expireDrop(int worldId, int channelId, int mapId, int id) {
      EventProducerRegistry.getInstance().send(DropExpiredEvent.class, EventConstants.TOPIC_DROP_EXPIRE_EVENT, worldId, channelId,
            new DropExpiredEvent(worldId, channelId, mapId, id));
   }
}
