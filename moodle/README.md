# Question types
## Drag and drop into text
This question type is implemented in the `DropText` type.

![Moodle rendering a 'Drag and drop into text' question](exampleImages/dropText.png)

## Drag and drop markers
This question type is implemented in the `DropMarker` type.

These questions contain an image onto which the markers are to be dropped. When creating one of these questions, the image must be base64 encoded. It is then bundled into the generated XML data.

## Multiple choice
This question type is implemented in the `MultiChoice` type. If there are more than a single correct answer, the package will automatically switch between 'One answer only' and 'Multiple answers allowed'.

![Moodle rendering a 'Multiple choice' question](exampleImages/multichoice.png)

## Numerical
This question type is implemented in the `Numerical` type.

## Short Answer
This question is implemented in the `ShortText` type.

![Moodle rendering a 'Short Answer' question](exampleImages/shortText.png)

For this type of question, it would be common to add specific feedback for partially correct answers. This is possible via the `NewAnswerWithFeedback` function.

![Specific feedback for partially correct answer](exampleImages/shortTextFeedback.png)
