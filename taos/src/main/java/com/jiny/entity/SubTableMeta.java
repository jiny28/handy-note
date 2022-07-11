package com.jiny.entity;

import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.util.List;

@Data
@AllArgsConstructor
@NoArgsConstructor
public class SubTableMeta {

    private String database;
    private String superTable;
    private String name;
    private List<TagValue> tags;
    private List<FieldMeta> fields;
}
