/* vim: set sw=8 ts=8 sts=8 noet: */
#include <capn.h>
#include <stdlib.h>

struct str {
	char *str;
	int len, cap;
};

extern char str_static[];
#define STR_INIT {str_static, 0, 0}

void str_reserve(struct str *v, int sz);

static void str_release(struct str *v) {
	if (v->cap) {
		free(v->str);
	}
}

static void str_reset(struct str *v) {
	if (v->len) {
		v->len = 0;
		v->str[0] = '\0';
	}
}

static void str_setlen(struct str *v, int sz) {
	str_reserve(v, sz);
	v->str[sz] = '\0';
	v->len = sz;
}

#ifdef __GNUC__
#define ATTR(FMT, ARGS) __attribute__((format(printf,FMT,ARGS)))
#else
#define ATTR(FMT, ARGS)
#endif

void str_add(struct str *v, const char *str, int sz);
int str_vaddf(struct str *v, const char *format, va_list ap) ATTR(2,0);
int str_addf(struct str *v, const char *format, ...) ATTR(2,3);
char *strf(struct str *v, const char *format, ...) ATTR(2,3);


