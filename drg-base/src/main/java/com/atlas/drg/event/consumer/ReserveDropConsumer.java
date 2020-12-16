package com.atlas.drg.event.consumer;

import com.atlas.drg.command.ReserveDropCommand;
import com.atlas.drg.constant.EventConstants;
import com.atlas.drg.processor.DropProcessor;
import com.atlas.drg.processor.TopicDiscoveryProcessor;
import com.atlas.kafka.consumer.SimpleEventHandler;

public class ReserveDropConsumer implements SimpleEventHandler<ReserveDropCommand> {
   @Override
   public void handle(Long key, ReserveDropCommand event) {
      DropProcessor.reserveDrop(event.dropId(), event.characterId());
   }

   @Override
   public Class<ReserveDropCommand> getEventClass() {
      return ReserveDropCommand.class;
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
      return TopicDiscoveryProcessor.getTopic(EventConstants.TOPIC_RESERVE_DROP_COMMAND);
   }
}
