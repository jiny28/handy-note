package com.jiny.entity;

import lombok.Data;

import java.util.List;

@Data
public class SubTableValue {

    private String database;
    private String superTable;
    private String name;
    private List<TagValue> tags;
    private List<RowValue> values;
}
