package moodle_test

import (
	"os"

	"github.com/ReneBoedker/MoodlishInquisition/moodle"
)

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
	// 	<shuffleanswers/>
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
