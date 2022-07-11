package com.jiny;

import com.jiny.api.SuperTableInterface;
import com.jiny.entity.FieldMeta;
import com.jiny.entity.SuperTableMeta;
import com.jiny.entity.TagMeta;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.context.SpringBootTest;

import java.util.ArrayList;

@SpringBootTest
class SuperTableTest {

    @Autowired
    private SuperTableInterface superTableInterface;



    @Test
    void createStable() {
        //CREATE STABLE meters (ts timestamp, current float, voltage int, phase float) TAGS (location binary(64), groupId int);
        superTableInterface.create(new SuperTableMeta("", "meters", new ArrayList<FieldMeta>() {{
            add(new FieldMeta("ts", "TIMESTAMP"));
            add(new FieldMeta("current", "FLOAT"));
            add(new FieldMeta("voltage", "INT"));
            add(new FieldMeta("phase", "FLOAT"));
        }}, new ArrayList<TagMeta>() {{
            add(new TagMeta("location", "BINARY(64)"));
            add(new TagMeta("groupId", "INT"));
        }}));

    }


    @Test
    void addFieldSTable() {
        superTableInterface.addField("", "meters", new FieldMeta("testf", "INT"));
        showSuperTable();
    }

    @Test
    void showSuperTable() {
        SuperTableMeta meters = superTableInterface.show("", "meters");
        System.out.println(meters);
    }

    @Test
    void delFieldSTable() {
        superTableInterface.delField("", "meters", "testf");
        showSuperTable();
    }


    @Test
    void addTagSTable() {
        superTableInterface.addTag("", "meters", new TagMeta("testf", "INT"));
        showSuperTable();
    }


    @Test
    void updateTagSTable() {
        superTableInterface.updateTag("", "meters", "testf", "tttt");
        showSuperTable();
    }

    @Test
    void delTagSTable() {
        superTableInterface.delTag("", "meters", "tttt");
        showSuperTable();
    }

    @Test
    void dropSTable() {
        superTableInterface.drop("", "meters");
        showSuperTable();
    }

    @Test
    void updateSubTableTag() {
        // 客户端需要配置FQDN，host设置映射到容器id或者是容器hostname
        //taos.updateSubTableTag("d1001", "groupid", "444");
        //
    }
}