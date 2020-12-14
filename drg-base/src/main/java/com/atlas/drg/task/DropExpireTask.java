package com.atlas.drg.task;

import com.atlas.drg.ConfigurationRegistry;
import com.atlas.drg.DropRegistry;
import com.atlas.drg.processor.DropProcessor;

public class DropExpireTask implements Runnable {
   @Override
   public void run() {
      int itemExpireInterval = ConfigurationRegistry.getInstance().getConfiguration().itemExpireInterval;

      DropRegistry.getInstance().getDrops().parallelStream()
            .filter(drop -> drop.dropTime() + itemExpireInterval < System.currentTimeMillis())
            .forEach(DropProcessor::destroyDrop);
   }
}
