package org.flowdev.flowparser.data;

import java.util.List;


public class Flow {
    private String name;
    private List<Operation> operations;
    private List<Connection> connections;

    public Flow name(final String name) {
        this.name = name;
        return this;
    }

    public Flow operations(final List<Operation> operations) {
        this.operations = operations;
        return this;
    }

    public Flow connections(final List<Connection> connections) {
        this.connections = connections;
        return this;
    }

    public String name() {
        return name;
    }

    public List<Operation> operations() {
        return operations;
    }

    public List<Connection> connections() {
        return connections;
    }
}
