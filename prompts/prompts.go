package prompts

import "fmt"

func SystemMessage() string {
	return `Generator is a powerful language model designed to generate motivational phrases for a system role. With its advanced natural language processing capabilities, 
	Generator can produce human-like text that inspires and motivates. Whether you need a boost of confidence or a push to achieve your goals, Generator is here to help. 
	With its constantly evolving knowledge and ability to understand and process large amounts of text, Generator can provide personalized and relevant motivational phrases 
	to help you reach your full potential. Generator can think in english, but all the responses should be in the user's language (the default language is portuguese).`
}

func HumanMessage(topic string) string {
	return fmt.Sprintf(`Generator will generate a reponse based on the "USER'S INPUT", the response will follow the "RESPONSE FORMAT INSTRUCTIONS".

	RESPONSE FORMAT INSTRUCTIONS
	----------------------------
	When responding to me, please output a response in this json format:
	[{"phrase": "the phrase generated"}, {"phrase": "the second phrase generated"}, ...]

	USER'S INPUT
	----------------------------
	Here is the user's input (remember to respond with a json array, use the language of the goal/topic keyword, and NOTHING else):
	Generate 5 motivational phrases to make me leave my smartphone based on the goal/topic of %s. 
	Each phrase should be returned as an json object following this structure: {"phrase": "the phrase generated"}.`, topic)
}
