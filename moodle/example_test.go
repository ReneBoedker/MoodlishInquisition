package moodle_test

import (
	"fmt"
	"os"

	"github.com/ReneBoedker/MoodlishInquisition/moodle"
)

// This simple example shows how to create a question bank containing a single
// question. The question has one correct answer (giving 100% of the available
// points).
func Example() {
	mc := moodle.NewMultiChoice(
		"What is your quest?",
		1,
		[]*moodle.Answer{
			moodle.NewAnswer("To seek the Holy Grail", 100),
			moodle.NewAnswer("Blue, no yel...", 0),
			moodle.NewAnswer("An African or European swallow", 0),
		},
	)

	qb := moodle.NewQuestionBank(
		"Bridge of Death",
		[]moodle.Question{mc},
	)

	qb.ToXml(os.Stdout)
	// Output:
	// <?xml version="1.0" encoding="UTF-8"?>
	// <quiz>
	// <question type="category">
	// 	<category>
	// 		<text>$module$/Bridge of Death</text>
	// 	</category>
	// </question>
	// <question type="multichoice">
	// 	<name>
	// 		<text>6C49FC87</text>
	// 	</name>
	// 	<questiontext format="html">
	// 		<text><![CDATA[What is your quest?]]></text>
	// 	</questiontext>
	// 	<defaultgrade>1</defaultgrade>
	// 	<shuffleanswers>1</shuffleanswers>
	// 	<answer fraction="100.000000">
	// 		<text><![CDATA[To seek the Holy Grail]]></text>
	// 	</answer>
	// 	<answer fraction="0.000000">
	// 		<text><![CDATA[Blue, no yel...]]></text>
	// 	</answer>
	// 	<answer fraction="0.000000">
	// 		<text><![CDATA[An African or European swallow]]></text>
	// 	</answer>
	// <single>true</single>
	// <answernumbering>none</answernumbering>
	// </question>
	// </quiz>
}

func ExampleEscapeMath() {
	fmt.Println(moodle.EscapeMath("$$x<5$$"))
	// Output:
	// $$x\lt 5$$
}

func ExampleNewDropText() {
	description := `In Latin, &ldquo;Romans go home&rdquo; is [[1]] [[3]] [[5]]`
	markers := []*moodle.TextMark{
		moodle.NewTextMark("Romani", 0, false),
		moodle.NewTextMark("Romanes", 0, false),
		moodle.NewTextMark("ite", 1, false),
		moodle.NewTextMark("eunt", 1, false),
		moodle.NewTextMark("domum", 2, false),
		moodle.NewTextMark("domus", 2, false),
	}

	question := moodle.NewDropText(description, 1, markers)

	question.ToXml(os.Stdout)
	// Output:
	// 	<question type="ddwtos">
	// 	<name>
	// 		<text>F56B8CB8</text>
	// 	</name>
	// 	<questiontext format="html">
	// 		<text><![CDATA[In Latin, &ldquo;Romans go home&rdquo; is [[1]] [[3]] [[5]]]]></text>
	// 	</questiontext>
	// 	<defaultgrade>1</defaultgrade>
	// 	<shuffleanswers>1</shuffleanswers>
	// 	<dragbox>
	// 		<text>Romani</text>
	// 		<group>1</group>
	// 	</dragbox>
	// 	<dragbox>
	// 		<text>Romanes</text>
	// 		<group>1</group>
	// 	</dragbox>
	// 	<dragbox>
	// 		<text>ite</text>
	// 		<group>2</group>
	// 	</dragbox>
	// 	<dragbox>
	// 		<text>eunt</text>
	// 		<group>2</group>
	// 	</dragbox>
	// 	<dragbox>
	// 		<text>domum</text>
	// 		<group>3</group>
	// 	</dragbox>
	// 	<dragbox>
	// 		<text>domus</text>
	// 		<group>3</group>
	// 	</dragbox>
	// </question>
}
