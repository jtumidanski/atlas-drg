package com.atlas.drg.rest;

import javax.ws.rs.Consumes;
import javax.ws.rs.GET;
import javax.ws.rs.Path;
import javax.ws.rs.PathParam;
import javax.ws.rs.Produces;
import javax.ws.rs.core.MediaType;
import javax.ws.rs.core.Response;

import com.atlas.drg.rest.processor.DropProcessor;

@Path("worlds")
public class WorldResource {
   @GET
   @Path("/{worldId}/channels/{channelId}/maps/{mapId}/drops")
   @Consumes(MediaType.APPLICATION_JSON)
   @Produces(MediaType.APPLICATION_JSON)
   public Response getDropsInMap(@PathParam("worldId") Integer worldId, @PathParam("channelId") Integer channelId,
                                 @PathParam("mapId") Integer mapId) {
      return DropProcessor.getDropsInMap(worldId, channelId, mapId).build();
   }
}
