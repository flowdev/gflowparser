package org.flowdev.flowparser.data;

public class PortData {
    private String name;
    private String capName;
    private boolean hasIndex;
    private int index;
    private int srcPos;


    public String name() {
        return this.name;
    }

    public PortData name(String name) {
        this.name = name;
        return this;
    }

    public String capName() {
        return this.capName;
    }

    public PortData capName(final String capName) {
        this.capName = capName;
        return this;
    }

    public boolean hasIndex() {
        return this.hasIndex;
    }

    public PortData hasIndex(boolean hasIndex) {
        this.hasIndex = hasIndex;
        return this;
    }

    public int index() {
        return this.index;
    }

    public PortData index(int index) {
        this.index = index;
        return this;
    }

    public int srcPos() {
        return this.srcPos;
    }

    public PortData srcPos(final int srcPos) {
        this.srcPos = srcPos;
        return this;
    }

}
