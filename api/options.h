typedef void * Options;
extern Options startOptions(char * config, char * user, char * overrides);
extern void addIntOption(Options, char *, int );
extern void addBoolOption(Options, char *, int);
extern void endOptions(Options );
