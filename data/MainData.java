package org.flowdev.flowparser.data;

import org.flowdev.parser.data.ParserData;

public class MainData {
    private ParserData parserData = new ParserData();
    private FlowFile flowFile;
    private String format;
    private String outputName;
    private String outputContent;

    public ParserData parserData() {
        return this.parserData;
    }

    public FlowFile flowFile() {
        return this.flowFile;
    }

    public String format() {
        return this.format;
    }

    public String outputContent() {
        return outputContent;
    }

    public String outputName() {
        return outputName;
    }

    public MainData parserData(ParserData parserData) {
        this.parserData = parserData;
        return this;
    }

    public MainData flowFile(FlowFile flowFile) {
        this.flowFile = flowFile;
        return this;
    }

    public MainData format(String format) {
        this.format = format;
        return this;
    }

    public MainData outputContent(String fileContent) {
        this.outputContent = fileContent;
        return this;
    }

    public MainData outputName(String fileName) {
        this.outputName = fileName;
        return this;
    }
}
