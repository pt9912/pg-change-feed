# mask-quotes.awk — Segmentierungs-Vorstufe des PreToolUse-Guards
# (.claude/hooks/pretooluse-command-guard.sh). Liest den Befehlsstring aus der
# Umgebungsvariablen CMD und ersetzt jeden Trenner (| & ; ( `), jeden Leerraum
# und jeden Zeilenumbruch INNERHALB von Anfuehrungszeichen oder hinter einem
# Backslash durch ein Steuerzeichen (\001 | \002 & \003 ; \004 ( \005 `
# \006 Leerraum \007 Zeilenumbruch); der Guard trennt danach nur noch an
# Trennern ausserhalb von Anfuehrungszeichen und liest ein Anfuehrungszeichen-
# Argument als ein Token. In doppelten Anfuehrungszeichen bleiben `$(` und
# Backtick lebendig (sie fuehren aus). `\;` bleibt ein Trenner (find-Ende).
# Ein unbalanciertes Anfuehrungszeichen ist Parse-Zweifel: die Eingabe geht
# unveraendert zurueck, der Guard segmentiert dann ohne Anfuehrungszeichen-
# Kenntnis (mehr Segmente, nie weniger).
function mk(c) { return (c in M) ? M[c] : c }
BEGIN {
  M["|"] = sprintf("%c", 1); M["&"] = sprintf("%c", 2); M[";"] = sprintf("%c", 3)
  M["("] = sprintf("%c", 4); M["`"] = sprintf("%c", 5); M[" "] = sprintf("%c", 6)
  M["\t"] = M[" "]; M["\n"] = sprintf("%c", 7); M["\r"] = M["\n"]
  s = ENVIRON["CMD"]; n = length(s); sp = 0; st[0] = "N"; out = ""
  for (i = 1; i <= n; i++) {
    c = substr(s, i, 1); t = st[sp]
    if (t == "S") { if (c == "'") sp--; else c = mk(c) }
    else if (c == "\\") {
      d = substr(s, ++i, 1)
      if (t != "D" && d == ";") c = "\\;"
      else if (t != "D" && d == "\n") c = " "
      else c = "\\" mk(d)
    }
    else if (t == "D") {
      if (c == "\"") sp--
      else if (c == "$" && substr(s, i + 1, 1) == "(") { st[++sp] = "P"; c = "$("; i++ }
      else if (c == "`") st[++sp] = "B"
      else c = mk(c)
    }
    else if (c == "'") st[++sp] = "S"
    else if (c == "\"") st[++sp] = "D"
    else if ((c == ")" && t == "P") || (c == "`" && t == "B")) sp--
    out = out c
  }
  printf "%s", (sp == 0 ? out : s)
}
