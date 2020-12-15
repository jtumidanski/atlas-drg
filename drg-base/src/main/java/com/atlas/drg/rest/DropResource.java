package com.atlas.drg.rest;

import javax.ws.rs.Consumes;
import javax.ws.rs.GET;
import javax.ws.rs.Path;
import javax.ws.rs.PathParam;
import javax.ws.rs.Produces;
import javax.ws.rs.core.MediaType;
import javax.ws.rs.core.Response;

import com.atlas.drg.rest.processor.DropProcessor;

@Path("drops")
public class DropResource {
   @GET
   @Path("/{id}}")
   @Consumes(MediaType.APPLICATION_JSON)
   @Produces(MediaType.APPLICATION_JSON)
   public Response getDropById(@PathParam("id") Integer id) {
      return DropProcessor.getDropById(id).build();
   }
}
