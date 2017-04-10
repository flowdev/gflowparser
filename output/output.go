package output

import (
	"github.com/flowdev/gflowparser/data"
)

var AllowedFormat = map[string]bool{
	"dot": true,
	"go":  true,
}
var specialFormat = map[string]bool{
	"dot": true,
}

// ------------ ProduceFormats:
// input:  *data.MainData{SelectedFormats, FlowFile}
// output: *data.MainData{SelectedFormats, FlowFile, CurrentFormat} to the correct output port.
type ProduceFormats struct {
	outPort        func(*data.MainData)
	specialOutPort func(*data.MainData)
}

func NewOutputFormats() *ProduceFormats {
	return &ProduceFormats{}
}
func (op *ProduceFormats) InPort(md *data.MainData) {
	for _, format := range md.SelectedFormats {
		md.CurrentFormat = format
		if specialFormat[format] {
			op.specialOutPort(md)
		} else {
			op.outPort(md)
		}
	}
}
func (op *ProduceFormats) SetOutPort(port func(*data.MainData)) {
	op.outPort = port
}
func (op *ProduceFormats) SetSpecialOutPort(port func(*data.MainData)) {
	op.specialOutPort = port
}

// ------------ FillPortPairs:
// input:  flows with operations with ports but without port pairs
// output: flows with operations with port pairs filled
type FillPortPairs struct {
	outPort func(*data.MainData)
}

func NewFillPortPairs() *FillPortPairs {
	return &FillPortPairs{}
}
func (op *FillPortPairs) InPort(md *data.MainData) {
	for _, flow := range md.FlowFile.Flows {
		for _, op := range flow.Ops {
			fillPortPairs4Op(op)
		}
	}
	op.outPort(md)
}
func (op *FillPortPairs) SetOutPort(port func(*data.MainData)) {
	op.outPort = port
}

func fillPortPairs4Op(op *data.Operation) {
	l := len(op.InPorts)
	m := len(op.OutPorts)
	n := max(l, m)
	portPairs := make([]*data.PortPair, n)
	for i := 0; i < n; i++ {
		p := &data.PortPair{}
		if i < l {
			p.InPort = op.InPorts[i]
		}
		if i < m {
			p.OutPort = op.OutPorts[i]
		}
		portPairs[i] = p
	}
}
func max(a, b int) int {
	if a >= b {
		return a
	}
	return b
}

// ------------ FillTemplate:
// input:  flows with operations with ports but without port pairs
// output: flows with operations with port pairs filled
type FillTemplate struct {
	outPort func(*data.MainData)
	tmplCache map[string]*template.
}

func NewFillTemplate() *FillTemplate {
	return &FillTemplate{}
}
func (op *FillTemplate) InPort(md *data.MainData) {
	for _, flow := range md.FlowFile.Flows {
		for _, op := range flow.Ops {
			fillPortPairs4Op(op)
		}
	}
	op.outPort(md)
}
func (op *FillTemplate) SetOutPort(port func(*data.MainData)) {
	op.outPort = port
}


package org.flowdev.flowparser.output;

import com.github.mustachejava.DefaultMustacheFactory;
import com.github.mustachejava.Mustache;
import com.github.mustachejava.MustacheFactory;
import org.flowdev.base.op.FilterOp;
import org.flowdev.flowparser.data.Flow;
import org.flowdev.flowparser.data.MainData;

import java.io.StringWriter;
import java.util.HashMap;
import java.util.Map;


public class FillTemplate extends FilterOp<MainData, FillTemplate.FillTemplateConfig> {
    private static final String TEMPLATE_DIR = FillTemplate.class.getPackage().getName().replace('.', '/') + "/";

    @Override
    protected void filter(MainData data) {
        MustacheFactory mf = new DefaultMustacheFactory(
                TEMPLATE_DIR + data.format());
        Mustache standardTpl = mf.compile("template.mustache");
        StringBuilder fileContent = new StringBuilder(8192);
        Map<String, Object> tplData = new HashMap<>();
        tplData.put("horizontal", getVolatileConfig().horizontal());

        for (Flow flow : data.flowFile().flows()) {
            tplData.put("flow", flow);
            StringWriter sw = new StringWriter();
            standardTpl.execute(sw, tplData);
            sw.flush();
            fileContent.append(sw.toString());
        }

        data.outputContent(fileContent.toString());
        outPort.send(data);
    }

    public static class FillTemplateConfig {
        private boolean horizontal;

        public boolean horizontal() {
            return horizontal;
        }

        public FillTemplateConfig horizontal(boolean horizontal) {
            this.horizontal = horizontal;
            return this;
        }
    }
}
