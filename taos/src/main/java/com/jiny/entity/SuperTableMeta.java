package com.jiny.entity;

import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.util.List;

@Data
@AllArgsConstructor
@NoArgsConstructor
public class SuperTableMeta {

    private String database;
    private String name;
    private List<FieldMeta> fields;
    private List<TagMeta> tags;
}