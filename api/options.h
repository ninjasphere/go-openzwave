extern void startOptions(char * config, char * user, char * overrides);
extern void addIntOption(char *, int );
extern void addBoolOption(char *, bool flag);
extern void addStringOption(char *, char * value, bool append);
extern void endOptions();
