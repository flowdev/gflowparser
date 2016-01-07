package org.flowdev.flowparser.data;

import java.util.ArrayList;
import java.util.List;


public class Operation {
    private String name;
    private String type;
    private List<PortData> inPorts = new ArrayList<>();
    private List<PortData> outPorts = new ArrayList<>();
    private int srcPos;
    private List<PortPair> portPairs;


    public Operation name(final String name) {
        this.name = name;
        return this;
    }

    public Operation type(final String type) {
        this.type = type;
        return this;
    }

    public String name() {
        return name;
    }

    public String type() {
        return type;
    }

    public List<PortData> inPorts() {
        return this.inPorts;
    }

    public Operation inPorts(final List<PortData> inPorts) {
        this.inPorts = inPorts;
        return this;
    }

    public List<PortData> outPorts() {
        return this.outPorts;
    }

    public Operation outPorts(final List<PortData> outPorts) {
        this.outPorts = outPorts;
        return this;
    }

    public int srcPos() {
        return this.srcPos;
    }

    public Operation srcPos(final int srcPos) {
        this.srcPos = srcPos;
        return this;
    }

    public List<PortPair> portPairs() {
        return portPairs;
    }

    public Operation portPairs(final List<PortPair> ports) {
        this.portPairs = ports;
        return this;
    }

}
