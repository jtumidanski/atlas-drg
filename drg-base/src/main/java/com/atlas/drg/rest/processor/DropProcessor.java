package com.atlas.drg.rest.processor;

import com.app.rest.util.stream.Collectors;
import com.atlas.drg.DropRegistry;

import builder.ResultBuilder;

public final class DropProcessor {
   private DropProcessor() {
   }

   public static ResultBuilder getDropsInMap(int worldId, int channelId, int mapId) {
      return DropRegistry.getInstance().getDropsForMap(worldId, channelId, mapId).stream()
            .map(ResultObjectFactory::create)
            .collect(Collectors.toResultBuilder());
   }
}
