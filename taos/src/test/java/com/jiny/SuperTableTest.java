package com.jiny;

import com.jiny.api.SuperTableInterface;
import com.jiny.entity.FieldMeta;
import com.jiny.entity.SuperTableMeta;
import com.jiny.entity.TagMeta;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.context.SpringBootTest;

import java.util.ArrayList;
import java.util.List;

@SpringBootTest
class SuperTableTest {

    @Autowired
    private SuperTableInterface superTableInterface;



    @Test
    void createStable() {
        //CREATE STABLE meters (ts timestamp, current float, voltage int, phase float) TAGS (location binary(64), groupId int);
        // generate fields
        for (int j = 0; j < 1000; j++) {
            List<FieldMeta> fieldMetas = new ArrayList<>();
            fieldMetas.add(new FieldMeta("ts", "TIMESTAMP"));
            for (int i = 0; i < 200; i++) {
                String fieldName = "item" + i;
                //String fieldType =i%2 == 0 ? "DOUBLE" : "BINARY(64)";
                String fieldType ="DOUBLE";
                fieldMetas.add(new FieldMeta(fieldName, fieldType));
            }
            SuperTableMeta superTableMeta = new SuperTableMeta("hlhz1", "device"+j, fieldMetas, new ArrayList<TagMeta>() {{
                add(new TagMeta("device", "BINARY(64)"));
            }});
            superTableInterface.create(superTableMeta);
        }
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
