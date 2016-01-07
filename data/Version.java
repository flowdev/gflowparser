package org.flowdev.flowparser.data;


public class Version {
    private int political;
    private int major;

    public int political() {
        return political;
    }

    public Version political(int political) {
        this.political = political;
        return this;
    }

    public int major() {
        return major;
    }

    public Version major(int major) {
        this.major = major;
        return this;
    }
}
