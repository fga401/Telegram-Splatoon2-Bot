package help

const (
	textKeyHelp = `
*Commands*:
- stages: /help\_stages`
	textKeyHelpStageSchedules = `
*Usage*:
/stages \[<prim\_filter>] \[<sec\_filters>...]

*<prim_filter>* should be:
- *[lgr]+* shows 'League', 'Gachi (Ranked)' or 'Regular'.

*<sec_filters>* could be:
- *\d+* shows the following N stage(s).
- *[ztrc]+* shows 'Splat Zones', 'Tower Control', 'Rainmaker' and 'Clam Blitz'.
- *b(\d+)-(\d+)* shows stages between X to Y o'clock.

_Default Case:_
- If no filter provided, it will add default filters 'lgr 1'.
- If no primary filter provided, it will add primary filters 'lgr'.
- If no secondary filter provided, it will add secondary filters '2'.
`
)
