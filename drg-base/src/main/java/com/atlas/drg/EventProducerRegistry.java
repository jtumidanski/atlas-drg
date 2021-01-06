package com.atlas.drg;

import java.util.HashMap;
import java.util.Map;

import com.atlas.drg.processor.TopicDiscoveryProcessor;
import com.atlas.kafka.KafkaProducerFactory;
import org.apache.kafka.clients.producer.Producer;
import org.apache.kafka.clients.producer.ProducerRecord;

public class EventProducerRegistry {
   private static final Object lock = new Object();

   private static final Object producerLock = new Object();

   private static volatile EventProducerRegistry instance;

   private final Map<String, String> topicMap;

   private Producer<Long, ?> producer;

   public static EventProducerRegistry getInstance() {
      EventProducerRegistry result = instance;
      if (result == null) {
         synchronized (lock) {
            result = instance;
            if (result == null) {
               result = new EventProducerRegistry();
               instance = result;
            }
         }
      }
      return result;
   }

   private EventProducerRegistry() {
      topicMap = new HashMap<>();
   }

   public <T> void send(String topic, long key, T event) {
      ProducerRecord<Long, T> record = new ProducerRecord<>(getTopic(topic), key, event);
      Producer<Long, T> producer = getProducer();
      producer.send(record);
   }

   public <T> void sendAndFlush(String topic, long key, T event) {
      ProducerRecord<Long, T> record = new ProducerRecord<>(getTopic(topic), key, event);
      Producer<Long, T> producer = getProducer();
      producer.send(record);
      producer.flush();
   }

   public <T> void sendAndFlushAndClose(String topic, long key, T event) {
      ProducerRecord<Long, T> record = new ProducerRecord<>(getTopic(topic), key, event);
      Producer<Long, T> producer = getProducer();
      producer.send(record);
      producer.flush();
      producer.close();
   }

   protected String getTopic(String id) {
      if (!topicMap.containsKey(id)) {
         synchronized (topicMap) {
            if (!topicMap.containsKey(id)) {
               topicMap.put(id, TopicDiscoveryProcessor.getTopic(id));
            }
         }
      }
      return topicMap.get(id);
   }

   protected <T> Producer<Long, T> getProducer() {
      if (producer == null) {
         synchronized (producerLock) {
            if (producer == null) {
               producer = KafkaProducerFactory.createProducer(getProducerId(), getBootstrapServers());
            }
         }
      }
      return (Producer<Long, T>) producer;
   }

   protected String getProducerId() {
      return "Drop Registry Service";
   }

   protected String getBootstrapServers() {
      return System.getenv("BOOTSTRAP_SERVERS");
   }
}