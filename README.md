## WhatsApp Bot

#### How to start ?

1. Install `SQLite` in your Commputer
2. Generate OpenAI API key at https://platform.openai.com/account/api-keys
3. Create `.env` file in the root directory
4. And then add variable named `OPENAI_API_KEY` with your OpenAI API key as its value
5. Run `go get` to install required modules
6. Run `start.sh` (Make sure you assign execute permission)

#### Current Features :
- ChatGPT Integration. Send message to my number with **/ai** in the beginning text, so that chatgpt will send response. For Example: ``/ai what is github?``
- Help menu, just type ``/help``
- etc.

##### Other Features still in development...