package com.atlas.drg;

import java.io.FileNotFoundException;
import java.io.FileReader;
import java.io.IOException;

import com.atlas.drg.configuration.Configuration;
import com.esotericsoftware.yamlbeans.YamlReader;

public class ConfigurationRegistry {
   private static final Object lock = new Object();

   private static volatile ConfigurationRegistry instance;

   private final Configuration configuration;

   public static ConfigurationRegistry getInstance() {
      ConfigurationRegistry result = instance;
      if (result == null) {
         synchronized (lock) {
            result = instance;
            if (result == null) {
               result = new ConfigurationRegistry();
               instance = result;
            }
         }
      }
      return result;
   }

   private ConfigurationRegistry() {
      String fileName = "/service/config.yaml";
      String message;
      try {
         YamlReader reader = new YamlReader(new FileReader(fileName));
         configuration = reader.read(Configuration.class);
         reader.close();
      } catch (FileNotFoundException var3) {
         message = "Could not read config file " + fileName + ": " + var3.getMessage();
         throw new RuntimeException(message);
      } catch (IOException var4) {
         message = "Could not successfully parse config file " + fileName + ": " + var4.getMessage();
         throw new RuntimeException(message);
      }
   }

   public Configuration getConfiguration() {
      return configuration;
   }
}
