typedef void * Options;
extern Options createOptions(char *, char *);
extern void addIntOption(Options, char *, int );
extern void addBoolOption(Options, char *, int);
extern Manager lockOptions(Options );
