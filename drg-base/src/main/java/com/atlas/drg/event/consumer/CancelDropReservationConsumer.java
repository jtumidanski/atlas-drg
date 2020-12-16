package com.atlas.drg.event.consumer;

import com.atlas.drg.command.CancelDropReservationCommand;
import com.atlas.drg.constant.EventConstants;
import com.atlas.drg.processor.DropProcessor;
import com.atlas.drg.processor.TopicDiscoveryProcessor;
import com.atlas.kafka.consumer.SimpleEventHandler;

public class CancelDropReservationConsumer implements SimpleEventHandler<CancelDropReservationCommand> {
   @Override
   public void handle(Long key, CancelDropReservationCommand command) {
      DropProcessor.cancelDropReservation(command.dropId(), command.characterId());
   }

   @Override
   public Class<CancelDropReservationCommand> getEventClass() {
      return CancelDropReservationCommand.class;
   }

   @Override
   public String getConsumerId() {
      return "Drop Registry";
   }

   @Override
   public String getBootstrapServers() {
      return System.getenv("BOOTSTRAP_SERVERS");
   }

   @Override
   public String getTopic() {
      return TopicDiscoveryProcessor.getTopic(EventConstants.TOPIC_CANCEL_DROP_RESERVATION_COMMAND);
   }
}
