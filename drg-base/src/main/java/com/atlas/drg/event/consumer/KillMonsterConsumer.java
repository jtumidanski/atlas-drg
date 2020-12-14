package com.atlas.drg.event.consumer;

import com.atlas.drg.processor.DropProcessor;
import com.atlas.kafka.consumer.SimpleEventHandler;
import com.atlas.morg.rest.constant.EventConstants;
import com.atlas.morg.rest.event.MonsterEvent;
import com.atlas.morg.rest.event.MonsterEventType;

public class KillMonsterConsumer implements SimpleEventHandler<MonsterEvent> {
   @Override
   public void handle(Long key, MonsterEvent event) {
      if (event.type().equals(MonsterEventType.KILLED)) {
         createDrops(event.uniqueId(), event.actorId());
      }
   }

   protected void createDrops(int uniqueId, int killerId) {
      DropProcessor.createDrops(uniqueId, killerId);
   }

   @Override
   public Class<MonsterEvent> getEventClass() {
      return MonsterEvent.class;
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
      return System.getenv(EventConstants.TOPIC_MONSTER_EVENT);
   }
}
