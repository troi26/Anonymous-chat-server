# Anonymous-chat-server
Programming with Go

Приложението се инсталира с "go get https://github.com/yoanpetrov97/Anonymous-chat-server"

Това е конзолно приложение, Анонимен чат сървър, в който потребителите се лог-ват с някакво потребителско име. 
За свързване със сървъра се използва команда "log-in <потребителско_име>", като ако името не се подаде изрично, се задава автоматично такова (User-1, User-2, и т.н.)

При успешно влизане в сървъра, потребителят трябва да отговори на определен брой въпроси, свързани с него.
За промяна на отговор се използва командата "change", при чието въвеждане, на екрана се показват възможности за промяна:
0 - смяна на всички отговори;
1 - име;
2 - възраст;
3 - местоживеене;
4 - пол;
5 - професия;
To be continued....

С командата "users" се извежда на екрана броя потребители в сървъра, както и колко от тях вмомента са в чат стая.

За свързване с произволна стая се използва командата "connect", която може да приема и аргумент - големина на стаята (максимален брой потребители).

За изход от чат се използва командата "bye" (Което се изпраща и към ответната страна), като се иска потвърждаване. 
Ако отсрещният потребител напусне чата, сървърът пита потребителя дали иска да се свърже с нов човек.

За изход от програмата се използва "exit".

В процеса на комуникиране между двама потребители, на всеки 50, 100, 250, 500, 1000 съобщения им се показват нови данни за другия човек, като с командата "expose" даден потребител може да изпрати веднага следващата подробност за себе си, а с "expose all" се показват всички подробности.
