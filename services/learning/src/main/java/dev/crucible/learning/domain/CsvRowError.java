package dev.crucible.learning.domain;

public record CsvRowError(
    int line,
    String field,
    String message
) {}
