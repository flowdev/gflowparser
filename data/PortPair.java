package org.flowdev.flowparser.data;


public class PortPair {
    private PortData inPort;
    private PortData outPort;
    private boolean isLast;

    public PortPair inPort(final PortData inPort) {
        this.inPort = inPort;
        return this;
    }

    public PortPair outPort(PortData outPort) {
        this.outPort = outPort;
        return this;
    }

    public PortPair isLast(final boolean isLast) {
        this.isLast = isLast;
        return this;
    }

    public PortData inPort() {
        return inPort;
    }

    public PortData outPort() {
        return outPort;
    }

    public boolean isLast() {
        return isLast;
    }
}
