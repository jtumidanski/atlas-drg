package com.atlas.drg.event;

public record DropExpiredEvent(int worldId, int channelId, int mapId, int uniqueId) {
}
