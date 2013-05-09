#include "str.h"
#include <stdio.h>
#include <stdarg.h>
#include <string.h>

#ifndef va_copy
#   ifdef _MSC_VER
#       define va_copy(d,s) d = s
#   elif defined __GNUC__
#       define va_copy(d,s) __builtin_va_copy(d,s)
#   else
#       error
#   endif
#endif

char str_static[] = "\0";

void str_reserve(struct str *v, int sz) {
	if (sz < v->cap)
		return;

	v->cap = (v->cap * 2) + 16;
	if (sz > v->cap) {
		v->cap = (sz + 8) & ~7;
	}

	if (v->str == str_static) {
		v->str = NULL;
	}

	v->str = realloc(v->str, v->cap + 1);
}

void str_add(struct str *v, const char *str, int sz) {
	if (sz < 0)
		sz = strlen(str);
	str_reserve(v, v->len + sz);
	memcpy(v->str+v->len, str, sz);
	v->len += sz;
	v->str[v->len] = '\0';
}

int str_vaddf(struct str *v, const char* format, va_list ap) {
	str_reserve(v, v->len + 1);

	for (;;) {
		int ret;

		char* buf = v->str + v->len;
		int bufsz = v->cap - v->len;

		va_list aq;
		va_copy(aq, ap);

		/* We initialise buf[bufsz] to \0 to detect when snprintf runs out of
		 * buffer by seeing whether it overwrites it.
		 */
		buf[bufsz] = '\0';
		ret = vsnprintf(buf, bufsz + 1, format, aq);

		if (ret > bufsz) {
			/* snprintf has told us the size of buffer required (ISO C99
			 * behavior)
			 */
			str_reserve(v, v->len + ret);

		} else if (ret >= 0) {
			/* success */
			v->len += ret;
			return ret;

		} else if (buf[bufsz] != '\0') {
			/* snprintf has returned an error but has written to the end of the
			 * buffer (MSVC behavior). The buffer is not large enough so grow
			 * and retry. This can also occur with a format error if it occurs
			 * right on the boundary, but then we grow the buffer and can
			 * figure out its an error next time around.
			 */
			str_reserve(v, v->len + bufsz + 1);

		} else {
			/* snprintf has returned an error but has not written to the last
			 * character in the buffer. We have a format error.
			 */
			return -1;
		}
	}
}


int str_addf(struct str *v, const char* format, ...) {
	va_list ap;
	va_start(ap, format);
	return str_vaddf(v, format, ap);
}

