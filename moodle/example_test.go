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
			moodle.NewAnswer("An African or European swallow?", 0),
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
	// 		<text>3078D5A8</text>
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
	// 		<text><![CDATA[An African or European swallow?]]></text>
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

func ExampleNewMultiChoice() {
	question := moodle.NewMultiChoice(
		"What is your quest?",
		1,
		[]*moodle.Answer{
			moodle.NewAnswer("To seek the Holy Grail", 100),
			moodle.NewAnswer("Blue, no yel...", 0),
			moodle.NewAnswer("An African or European swallow?", 0),
		},
	)

	question.ToXml(os.Stdout)
	// Output:
	// 	<question type="multichoice">
	// 	<name>
	// 		<text>3078D5A8</text>
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
	// 		<text><![CDATA[An African or European swallow?]]></text>
	// 	</answer>
	// <single>true</single>
	// <answernumbering>none</answernumbering>
	// </question>
}

func ExampleNewShortText() {
	question := moodle.NewShortText(
		"Name one of the fresh fruits any self-defense course should cover.",
		1,
		[]*moodle.Answer{
			moodle.NewAnswer("passion fruit", 100),
			moodle.NewAnswer("orange", 100),
			moodle.NewAnswer("apple", 100),
			moodle.NewAnswer("whole grapefruit", 100),
			moodle.NewAnswer("grapefruit segments", 100),
			moodle.NewAnswerWithFeedback("*grapefruit*", 50, "Which types of grapefruit?"),
			moodle.NewAnswer("pomegranate", 100),
			moodle.NewAnswer("greengage", 100),
			moodle.NewAnswer("grape", 100),
			moodle.NewAnswer("lemon", 100),
			moodle.NewAnswer("plum", 100),
			moodle.NewAnswer("mangos in syrup", 100),
			moodle.NewAnswer("red cherry", 100),
			moodle.NewAnswer("black cherry", 100),
			moodle.NewAnswerWithFeedback("*cherry*", 50, "Which type of cherry?"),
			moodle.NewAnswer("banana", 100),
		},
	)

	question.SetCaseSensitivity(false)

	question.ToXml(os.Stdout)
	// Output:
	// 	<question type="shortanswer">
	// 	<name>
	// 		<text>BA90058D</text>
	// 	</name>
	// 	<questiontext format="html">
	// 		<text><![CDATA[Name one of the fresh fruits any self-defense course should cover.]]></text>
	// 	</questiontext>
	// 	<defaultgrade>1</defaultgrade>
	// 	<answer fraction="100.000000">
	// 		<text><![CDATA[passion fruit]]></text>
	// 	</answer>
	// 	<answer fraction="100.000000">
	// 		<text><![CDATA[orange]]></text>
	// 	</answer>
	// 	<answer fraction="100.000000">
	// 		<text><![CDATA[apple]]></text>
	// 	</answer>
	// 	<answer fraction="100.000000">
	// 		<text><![CDATA[whole grapefruit]]></text>
	// 	</answer>
	// 	<answer fraction="100.000000">
	// 		<text><![CDATA[grapefruit segments]]></text>
	// 	</answer>
	// 	<answer fraction="50.000000">
	// 		<text><![CDATA[*grapefruit*]]></text>
	// 		<feedback format="html">
	// 			<text><![CDATA[Which types of grapefruit?]]></text>
	// 		</feedback>
	// 	</answer>
	// 	<answer fraction="100.000000">
	// 		<text><![CDATA[pomegranate]]></text>
	// 	</answer>
	// 	<answer fraction="100.000000">
	// 		<text><![CDATA[greengage]]></text>
	// 	</answer>
	// 	<answer fraction="100.000000">
	// 		<text><![CDATA[grape]]></text>
	// 	</answer>
	// 	<answer fraction="100.000000">
	// 		<text><![CDATA[lemon]]></text>
	// 	</answer>
	// 	<answer fraction="100.000000">
	// 		<text><![CDATA[plum]]></text>
	// 	</answer>
	// 	<answer fraction="100.000000">
	// 		<text><![CDATA[mangos in syrup]]></text>
	// 	</answer>
	// 	<answer fraction="100.000000">
	// 		<text><![CDATA[red cherry]]></text>
	// 	</answer>
	// 	<answer fraction="100.000000">
	// 		<text><![CDATA[black cherry]]></text>
	// 	</answer>
	// 	<answer fraction="50.000000">
	// 		<text><![CDATA[*cherry*]]></text>
	// 		<feedback format="html">
	// 			<text><![CDATA[Which type of cherry?]]></text>
	// 		</feedback>
	// 	</answer>
	// 	<answer fraction="100.000000">
	// 		<text><![CDATA[banana]]></text>
	// 	</answer>
	// <usecase>0</usecase>
	// </question>
}
