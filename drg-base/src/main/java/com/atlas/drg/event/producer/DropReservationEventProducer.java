package com.atlas.drg.event.producer;

import com.atlas.drg.EventProducerRegistry;
import com.atlas.drg.constant.EventConstants;
import com.atlas.drg.event.DropReservationEvent;
import com.atlas.drg.event.DropReservationType;

public final class DropReservationEventProducer {
   private DropReservationEventProducer() {
   }

   public static void reservationSuccess(int dropId, int characterId) {
      reservationEvent(dropId, characterId, DropReservationType.SUCCESS);
   }

   public static void reservationFailure(int dropId, int characterId) {
      reservationEvent(dropId, characterId, DropReservationType.FAILURE);
   }

   protected static void reservationEvent(int dropId, int characterId, DropReservationType type) {
      EventProducerRegistry.getInstance()
            .send(DropReservationEvent.class, EventConstants.TOPIC_DROP_RESERVATION_EVENT, dropId,
                  new DropReservationEvent(characterId, dropId, type));
   }
}
