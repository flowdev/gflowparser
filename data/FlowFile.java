package org.flowdev.flowparser.data;

import java.util.List;


public class FlowFile {
    private String fileName;
    private Version version;
    private List<Flow> flows;

    public FlowFile fileName(final String fileName) {
        this.fileName = fileName;
        return this;
    }

    public FlowFile version(final Version version) {
        this.version = version;
        return this;
    }

    public FlowFile flows(final List<Flow> flows) {
        this.flows = flows;
        return this;
    }

    public String fileName() {
        return fileName;
    }

    public Version version() {
        return version;
    }

    public List<Flow> flows() {
        return flows;
    }
}
