typedef void * Options;
extern Options startOptions(char *, char *);
extern void addIntOption(Options, char *, int );
extern void addBoolOption(Options, char *, int);
extern void endOptions(Options );
