package com.atlas.drg.rest.processor;

import javax.ws.rs.core.Response;

import com.app.rest.util.stream.Collectors;
import com.app.rest.util.stream.Mappers;
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

   public static ResultBuilder getDropById(int id) {
      return DropRegistry.getInstance().getDrop(id)
            .map(ResultObjectFactory::create)
            .map(Mappers::singleOkResult)
            .orElse(new ResultBuilder(Response.Status.NOT_FOUND));
   }
}
