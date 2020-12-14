package com.atlas.drg.event;

public record DropEvent(int worldId, int channelId, int mapId, int uniqueId, int itemId, int meso, int dropType,
                        int dropX, int dropY,
                        int ownerId, Integer ownerPartyId, long dropTime,
                        int dropperUniqueId, int dropperX, int dropperY,
                        boolean playerDrop, byte mod) {
}
