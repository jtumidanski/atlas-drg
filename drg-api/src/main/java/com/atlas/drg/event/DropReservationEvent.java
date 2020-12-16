package com.atlas.drg.event;

public record DropReservationEvent(int characterId, int dropId, DropReservationType type) {
}
