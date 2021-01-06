package com.atlas.drg.event.producer;

import com.atlas.drg.EventProducerRegistry;
import com.atlas.drg.constant.EventConstants;
import com.atlas.drg.event.DropPickedUpEvent;

public final class DropPickedUpEventProducer {
   private DropPickedUpEventProducer() {
   }

   public static void emit(int dropId, int characterId, int mapId) {
      EventProducerRegistry.getInstance().send(EventConstants.TOPIC_PICKUP_DROP_EVENT, dropId,
            new DropPickedUpEvent(dropId, characterId, mapId));
   }
}
