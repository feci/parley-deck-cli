# Design critique & selection methodology — research digest for `parley-design` / `parley-design-check`

**Author:** research subagent · **Date:** 2026-07-28
**Question:** is our proposed multi-agent method — DIVERGE (each agent a *different* direction) → adversarial cross-critique → ONE direction WINS WHOLE → graft 2–3 details from losers — grounded in real human design practice, or invented?
**Answer up front:** it is *almost exactly* the Google Ventures Design Sprint "Sticky Decision" ritual, plus the Pixar Braintrust's authority rule, plus the agency "three concepts" convention. Every one of the four moves has a named precedent. Two moves need correction (see §6.3 and §8).

Legend: **[F]** = fact verified against the cited source. **[I]** = my inference/synthesis, not in any source.
Secondary sources are marked `(secondary)` where I could not reach the primary (e.g. the *Sprint* book itself, or a 403-walled page).

---

# 1. Design critique methodology at real orgs

## 1.1 GV (Google Ventures) — Braden Kowitz, "Guide to Design Critique"

Source: <https://library.gv.com/guide-to-design-critique-86ebf499bed5> (fetched via Medium mirror)

**[F] Five-step agenda:**
1. "Use formal critiques to jumpstart change" — scheduled sessions, **5–6 diverse participants**
2. "Start with critique guidelines" — ground rules stated *before* feedback begins
3. "Set the stage" — align on business/customer goals, constraints, schedule, and **fidelity level**
4. "Simulate the customer experience" — show **complete task flows, not isolated mockups**
5. "Gather feedback and discuss" — **collect written feedback first, then discuss as a group**

**[F] Verbatim rules:**
- Candor: *"It doesn't help anyone to stay silent during a critique, only to express your doubts privately later."*
- Goals over taste: *"Critique not about whether you like or dislike a design. Good feedback is about how the design is meeting (or missing) the customer and business goals."*
- Problems before solutions: *"Instead of arguing over solutions, start by taking a step back and first discussing problems with the current design."*
- No mandates: *"Any new ideas should be given as suggestions, not mandates."*
- Not a design session: *"We are not here to… design better solutions on the spot. Critiques help designers improve their designs."*

**[I]** Step 5 ("written feedback first, then discuss") is the single most important structural detail for us: it is the human analogue of "each agent writes its own artifact before reading the others'". It exists specifically to stop the first speaker from anchoring the room.

## 1.2 Figma — Noah Levin, VP Product Design

Sources: <https://www.figma.com/blog/design-critiques-at-figma/> and <https://nlevin.medium.com/design-critiques-at-figma-799d4a3a1b0>

**[F] Six named critique formats:**

| # | Format | Verbatim mechanic |
|---|---|---|
| 1 | **Standard** | "Present context, receive feedback." Context 10–15 min → clarifying Qs 2–5 min → feedback 10+ min |
| 2 | **Jams / Workshops** | "Brainstorms, Crazy 8's, Mood boarding, Sketching, etc." — "30 minutes heads down" then group discussion |
| 3 | **Pair Design** | "Work in small groups of 2–3." |
| 4 | **Silent Critiques** | "Everyone stays silent and reviews a document and adds feedback digitally" — enables "multiple conversation threads simultaneously, allowing for a greater volume of feedback" |
| 5 | **Paper / print-out** | "Print out work on paper and hang it up around the room." (Physical wall + stickies; ≤ monthly, high overhead) |
| 6 | **FYI** | "Sharing quick context on a project. No feedback necessary." |

**[F] Two feedback-turn disciplines, chosen by the *presenter*:**
- **Round-The-Room (RTR 🎡):** "Go one at a time around the room… Try to keep it to two minutes per person."
- **Popcorn (🍿):** "Freeform discussion… comments pop up from all over the place unpredictably."

**[F] Load limits:** "We've recently started only allowing a maximum of **two topics per one hour** critique." A timer is used. A named note-taker summarises.
**[F] Ground rules (Figma blog):** focus on goals, not preferences; observe before judging; ask before suggesting; be specific.
**[F] Philosophy:** *"critiques should be motivating, not intimidating. Helpful, not discouraging."*

**[I]** Figma's **Silent Critique** is the closest human format to what Parley Deck actually is — asynchronous, written, parallel threads, no turn-taking. We should name our critique round after it rather than after the spoken "crit", because the spoken-crit literature is full of rules (RTR timeboxes, popcorn) that solve a problem (one mouth at a time) that headless agents do not have.

## 1.3 Pixar — the Braintrust, and "plussing"

Sources: <https://collider.com/pixar-braintrust-details-ed-catmull-inside-out/>, <https://shohawk.com/the-pixar-braintrust/>, <https://medium.com/great-business-stories/lessons-from-pixar-1-the-braintrust-e306843a5153> (secondary — the Fast Company excerpt of *Creativity, Inc.* returned HTTP 403)

**[F] The two structural rules that define the Braintrust:**
1. **It has no authority.** "The director does not have to follow any of the specific suggestions, and after a Braintrust meeting, it is up to him or her to figure out how to address the feedback."
2. **Notes, not prescriptions.** "The Braintrust does not prescribe solutions, fixes or ideas… how to fix issues [is] the responsibility of the project director and their creative team." Members are "experts with empathy, who focus on identifying and diagnosing problems, not prescribing solutions."

**[F]** Candor (not merely honesty) is the stated currency: by stripping the group of formal authority, candor becomes *safe*, and "the only currency being the quality of the observation, not the rank of the person making it."
**[F]** Membership: people with deep storytelling understanding who "have been through the process themselves."

**[F] "Plussing"** (Pixar, inherited from Walt Disney's production meetings): *"if you give a criticism, you also have to give a plus — a suggestion for how to fix the criticism you just pointed out."* Rooted in improv's "yes, and". Sources: <https://intenseminimalism.com/2015/pixars-plussing-technique-of-giving-feedback/>, <https://leadx.org/articles/pixar/>, <https://www.andycleff.com/2014/09/plussing-learning-and-working-in-a-collaborative-environment/>

**[I] ⚠️ These two Pixar rules directly contradict each other and most write-ups miss it.** "Notes not prescriptions" says *do not* supply the fix; "plussing" says you *may not criticise without* supplying a fix. They are reconcilable only if you separate them by **phase and by bindingness**:
- Braintrust/notes-not-prescriptions governs **diagnosis of someone else's owned work** (the author owns the fix).
- Plussing governs **live co-ideation** (a jam, where nothing is owned yet).
Our protocol must pick one per phase, explicitly, and must mark any offered fix as **non-binding** (which is exactly GV's "suggestions, not mandates"). Do not let an agent's proposed fix become an implicit requirement.

## 1.4 Shopify / Aaron Irizarry (NASDAQ; co-author *Discussing Design*) — 12 Rules of Engagement

Source: <https://www.shopify.com/partners/blog/12-rules-of-engagement-when-running-design-critiques-from-nasdaq-s-aaron-irizarry>

**[F] Giving critique (6):** (1) **Don't invite yourself** — "make sure that it's actually wanted"; (2) **Use a filter** — feedback must be "objective and rooted in critical thinking, not based on emotional or subjective impressions"; (3) **Avoid problem solving** — critique is "analyzing and evaluating work, rather than actual problem solving"; (4) **Don't make assumptions** — learn the constraints first; (5) **Lead with questions**; (6) **Talk about strengths** — "describe what is working".

**[F] Receiving critique (6):** (1) **Set the foundation** (ground rules + objectives + personas agreed beforehand); (2) **Remember the purpose** — "critiquing is about understanding and improvement, not judgement"; (3) **Participate** — name the specific areas you want critiqued; (4) **Think before responding**; (5) **Use active listening** — "repeat their advice back… to confirm that you've fully understood"; (6) **Accept insight from all angles**.

**[I]** #3-receiving ("name the specific areas you want critiqued") maps to a typed field on our artifact: `critique_requests: [...]`. Figma does the same thing (presenter picks RTR vs popcorn, and states what feedback they want). This is a cheap, high-value protocol field.

## 1.5 Google (internal UX critique) — Anna Iurchenko

Source: <https://medium.com/google-design/a-collaborative-approach-to-shaping-successful-ux-critique-practices-b7f060c21582>

**[F]** Three roles in a critique: **facilitator, presenter, reviewer**.
**[F]** A tenet: **"Never established"** — the format, tools and feedback strategies are themselves permanently subject to experiment.
**[F]** A single signup document holds "purpose, principles, rules, and notes" — a one-stop-shop artifact for the critique.
**[F]** Preparation questions: "What problem am I trying to solve?" and "What aspects of your designs are you seeking feedback on?"

## 1.6 Stanford d.school — "I Like / I Wish / What If"

Sources: <https://spin.atomicobject.com/2018/09/12/i-like-i-wish-what-if/>, <https://public-media.interaction-design.org/pdf/I-Like-I-Wish-What-If.pdf> (PDF not text-extractable), <https://blog.uxtweak.com/design-critique/>

**[F]** Feedback is given in **I-statements**: *I like…* (what worked), *I wish…* (change, stated diplomatically), *What if…* (speculative/generative, no obligation).
**[F]** Common facilitation: ~10 minutes for each participant to fill **2–5 Post-its per category**, then discuss one-by-one with uninterrupted airtime per person.

**[I]** The three stems are a **typed critique record with three severity/modality classes**: `like` = preserve-this (a graft candidate!), `wish` = defect, `what_if` = non-binding proposal. That is precisely the Braintrust/plussing reconciliation from §1.3, already typed. `I like` is the most under-used of the three for us: it is *literally the graft harvest mechanism*.

## 1.7 De Bono — Six Thinking Hats (the origin of "black hat critique")

Sources: <https://www.designorate.com/the-six-hats-of-critical-thinking-and-how-to-use-them/>, <https://blog.prototypr.io/six-thinking-hats-method-for-constructive-design-feedback-c6285bced9f3>, <https://umnlibraries.manifoldapp.org/read/effective-design-critique-strategies-across-disciplines-...>

**[F]** Six roles: **White** (facts), **Red** (emotions/intuition), **Black** (critical judgement, risks, worst case), **Yellow** (optimism, positive impact), **Green** (creative generation), **Blue** (process control / meta).
**[F]** It is *parallel thinking*: the whole group wears the same hat at the same time, "focusing on one perspective at a time".
**[F]** Key justification, verbatim: the roles are **"blatantly artificial, a feature which helps separate individual ego from the activity, making it clear that statements made while wearing the hats do not necessarily reflect that person's unrestrained opinions."**

**[I]** That last sentence is the single strongest argument for *assigning* critique lenses to agents rather than asking for "your honest opinion". An agent asked for its opinion will defend its own proposal; an agent assigned the Black hat on a rival proposal is discharging a role. It also gives us a principled reason to make the critique **adversarial by assignment**, not by temperament.

---

# 2. Parallel divergence — the practice, and the exact GV decision ritual

## 2.1 Design charrette (NN/g's concrete version)

Source: <https://www.nngroup.com/articles/design-charrettes/>

**[F] Procedure, verbatim numbers:**
1. Group size: "as many as 20 people and as few as 2"; "Anyone can (and should) participate, not just UX designers."
2. "Write a goal or a design challenge on the whiteboard."
3. "Each person sketches his or her own ideas for **5 minutes**… all pens down."
4. "Each person works alone. **No talking once sketching begins.**"
5. "Each person gets **2 minutes (and no more than 2 minutes)** to show his or her ideas."
6. Q&A: "spending **one more minute** on each person."
7. "The UX designer collects the papers and uses the ideas generated to help derive a design."
**[F]** Critical success factor: "The person running the meeting has to keep time and be diligent about it. Otherwise these meetings can go forever."

**[F] Origin** (<https://en.wikipedia.org/wiki/Charrette>): Paris, 1800s — students' exams were collected in a cart (*charrette*) and some kept sketching madly as the cart came round. Charrettes commonly run 1–3 days and split into sub-groups that present back to the full group (<https://participedia.net/method/charrette>).

**[I]** Note step 7: the charrette **does not vote**. A single named designer collects everything and derives the design. That is a second, weaker precedent for "one owner integrates" — but it is *integration*, not "one winner whole", so it is a precedent for the **graft** half of our rule, not the winner half.

## 2.2 GV Design Sprint — Tuesday's "Four-Step Sketch" (the DIVERGE ritual)

Sources: <https://www.sessionlab.com/methods/four-step-sketch>, <https://maa1.medium.com/book-review-sprint-part-3-day-2-720b2be7c9c9>, <https://library.gv.com/sprint-week-tuesday-d22b30f905c3> (secondary summaries of Knapp/Zeratsky/Kowitz, *Sprint*, 2016)

**[F] Four steps, with timeboxes:**
| Step | Timebox | What |
|---|---|---|
| 1. **Notes** | 20 min | Silently walk the room, gather notes from Monday's material |
| 2. **Ideas** | 20 min | Privately jot rough ideas, circle the most promising |
| 3. **Crazy 8s** | 8 min | Fold a sheet into 8 frames; **one variation per frame, 1 minute each**, all of *your own best idea* |
| 4. **Solution Sketch** | 30–90 min | A **three-panel storyboard**, **self-explanatory**, and **kept anonymous** |

**[F]** The stated rationale, verbatim: **"We know that individuals working alone generate better solutions than groups brainstorming out loud."** Everyone does Notes / Ideas / Crazy 8s **without having shared anything with the rest of the team**. This is GV's "work alone together".

**[I] Two details are gold for us and are cheap to copy:**
- **Anonymity of the solution sketch.** GV anonymises *before* voting. We can too (see §7 P2). It kills self-preference and reputation effects at zero cost.
- **Crazy 8s is intra-agent divergence, not inter-agent.** Each person diverges *from their own best idea* 8 ways before committing. Our current plan only has inter-agent divergence. Adding a cheap intra-agent Crazy-8s pass (8 one-line variations before committing to a direction) is a documented, timeboxed way to widen each agent's own search before it locks in.

## 2.3 GV Design Sprint — Wednesday's "Sticky Decision" (the DECIDE ritual) ★

This is the ritual the task asked me to describe exactly. Sources: <https://medium.com/design-bootcamp/how-the-sticky-decision-works-in-a-sprint-6a8c79578f7b>, <https://medium.com/@mariajmorales312/sprint-wednesday-sticky-decision-8530595f17d>, <https://www.mentorist.app/action/structure-your-decision-making_541/> (secondary summaries of *Sprint*)

**[F] Five sub-steps, in fixed order:**

**① Art Museum.** Tape all solution sketches to the wall **in one row, "just like the paintings in a museum"**, in chronological/storyboard order. No presentation. No author names.

**② Heat Map.** **20–30 dot stickers per person.** Everyone works **silently**, simultaneously. Put dots next to *parts* you like — **2–3 dots on the most exciting ideas**. Write concerns on sticky notes *below* the sketch. Move on. → Output: a wall where interest is visible **before anyone has spoken**.

**③ Speed Critique.** **3 minutes per sketch** (extendable if the sketch is dense). Roles: **Facilitator narrates**, **Scribe** captures ideas on stickies. Sequence:
   1. Facilitator describes the sketch aloud;
   2. Facilitator calls out the ideas that attracted **dot clusters**;
   3. Team names anything missed;
   4. Scribe labels each idea with a short handle ("Animated Video");
   5. **The sketch's creator remains silent until the end**;
   6. Only then does the creator explain what was missed and answer questions.

**④ Straw Poll.** **~10 minutes. One vote per person** (a pink dot). Each person **privately writes** their choice first (may be a whole sketch or a single idea inside one), then places the dot, then gets **~1 minute** to say why. → Explicitly **advisory**.

**⑤ Supervote.** The **Decider** gets **three special votes**, marked with their initials. The Decider **may ignore the straw poll entirely** and may distribute the three votes however they like. Sketches with a supervote win; the rest become **"maybe-laters"** (kept, not deleted).

**[F] The "Rumble or All-in-One" branch** (immediately after the supervote). Sources: <https://coda.io/@jazer/design-sprint/2-rumble-30>, <https://www.gv.com/sprint/>, <https://brandgenetics.com/human-thinking/sprint-by-google-ventures-speed-summary/>
- If two winning concepts **conflict**: build **both** and test head-to-head — a **Rumble**. "Your prototypes will battle head-to-head, like professional wrestlers whacking each other with folding chairs."
- Rumble prototypes get **fake brand names**, generated by Note-and-Vote, so testers can tell them apart — otherwise "you risk sounding like an optometrist: 'Which version do you prefer? A, or B?'"
- If the winners **combine** cleanly: **"If you think you can combine your winning sketches into one product, don't bother with a rumble. Instead, put them together into your best shot at solving the problem."** → **All-in-One**.

**[I]** ⚠️ **This is the correction to our stated method.** GV does *not* say "one winner whole, always". GV says: **conflicting → Rumble (both survive, decided by external evidence); compatible → All-in-One (merge).** Our "one direction WINS WHOLE, graft 2–3 details from losers" is exactly GV's **All-in-One** path. We are missing the **Rumble** path — the case where two directions are genuinely incommensurable and the honest answer is *build both cheaply and let evidence decide*, not *pick one by vote*. See §7, phase P6b.

## 2.4 GV — "Note-and-Vote" (the general-purpose micro-ritual, ~15 min)

Source: Jake Knapp, <https://www.linkedin.com/pulse/note-and-vote-how-avoid-groupthink-meetings-jake-knapp> (also Fast Company <https://www.fastcompany.com/3034772/note-and-vote-how-google-ventures-avoids-groupthink-in-meetings>)

**[F] Seven steps, verbatim:**
1. **Note (5–10 min)** — "Everyone writes down as many ideas as they can. **Individually. Quietly.**"
2. **Self-edit (2 min)** — "Each person reviews his or her own list and picks **one or two favorites**. Individually. Quietly."
3. **Share and Capture** — "One at a time, each person shares his or her top idea(s). **No sales pitch.**"
4. **Vote (5 min)** — "Each person chooses a favorite from the ideas on the whiteboard. Individually. Quietly. **You must commit your vote to paper.**"
5. **Share and Capture (Votes)** — "One at a time, each person says their vote… Say what you wrote."
6. **Decide** — "**Who is the decider? She should make the final call — not the group.**"
7. **Rejoice** — "That only took 15 minutes!"

**[F]** Stated purpose: "short-circuit the worst parts of groupthink while getting the most out of different perspectives"; "By having most of the thinking done individually, you'll see a big efficiency boost"; "The votes are for information to give advice to the decider **who remains in charge**."

**[I]** Two mechanics generalise perfectly to agents: **"No sales pitch"** (share ≠ advocate) and **"commit your vote to paper" before revealing** (simultaneous reveal kills bandwagoning). Both are one-line protocol rules.

## 2.5 The agency "three concepts" convention

**[F]** I could **not** find an authoritative, citable source establishing "three concepts" as a codified convention; the sources returned are agency marketing pages, not methodology. **Treat "always show three" as folklore, not evidence.**
**[F]** What *is* evidenced nearby: charrettes explicitly "divide into sub-groups, and each sub-group then presents its work to the full group as material for further dialogue" (<https://participedia.net/method/charrette>) — i.e. parallel *teams*, N determined by group size, not by a magic 3.
**[I]** Do not justify our N (number of divergent directions) by appealing to "agencies always do three". Justify it by roster size and by §4's evaluator-count evidence (Nielsen: **3–5 evaluators find ~75% of issues at 5**). N = number of roster agents, floor 3, is defensible on that basis.

---

# 3. Why design-by-committee fails — what is actually documented

## 3.1 The proverb (weak evidence — say so)

**[F]** "A camel is a horse designed by a committee." Earliest located match: ***Reader's Digest*, September 1954**; the leading candidate for the core coinage is **T. R. Quaife**; the "horse" was added later by an anonymous improver. Frequently **mis**attributed to Charles F. Kettering and to Alec Issigonis (the Mini's designer) — "**no solid citations exist**". Source: Quote Investigator, <https://quoteinvestigator.com/2023/09/22/horse-camel/>
**[I]** Cite it as *rhetoric with a known provenance*, never as evidence. If `parley-design` quotes it, it must quote the QI provenance too, or it becomes exactly the kind of confident-but-unsourced claim our own protocol §13 penalises.

## 3.2 The real evidence — group *generation* loses (production blocking)

**[F]** Diehl & Stroebe (1987), *Productivity Loss in Brainstorming Groups: Toward the Solution of a Riddle*, JPSP. Four experiments testing free riding, evaluation apprehension and **production blocking**. Experiment 4 manipulated blocking directly and found **production blocking accounted for most of the productivity loss** of real (interacting) groups.
Source PDF: <https://homepages.se.edu/cvonbergen/files/2013/01/Productivity-Loss-In-Brainstorming_Toward-the-Solution-of-a-Riddle.pdf>
**[F]** Meta-analytic effect sizes for **nominal (independent) groups beating face-to-face groups**, 1958–1990: **r = .57** for number of non-redundant ideas and **r = .56** for idea quality. (Reported in the brainstorming meta-analysis literature; see also <https://en.wikipedia.org/wiki/Production_blocking>.)
**[F]** Evaluation apprehension was **not** supported as a major cause (no interaction found in Experiment 3).

**[I]** This validates the *diverge-in-isolation* half of our method with a large, replicated effect. Note the mechanism though: production blocking is "can't generate while listening" — a **serial-channel** constraint that headless agents in separate processes **do not have**. So we cannot borrow Diehl & Stroebe's *mechanism*; we can only borrow its *conclusion*, and we need a different mechanism to justify isolation for LLMs. §5 supplies it: **diversity collapse / stance homogenization**, which *is* an LLM-native mechanism and points the same way.

## 3.3 The real evidence — group *selection* loses (feasibility bias) ★

**[F]** Rietzschel, Nijstad & Stroebe (2010), *The selection of creative ideas after individual idea generation: choosing between creativity and impact*, **British Journal of Psychology**. PubMed: <https://pubmed.ncbi.nlm.nih.gov/19267959/> · Wiley: <https://bpspsychub.onlinelibrary.wiley.com/doi/10.1348/000712609X414204>
Finding: participants show a **"strong tendency to select feasible and desirable ideas, at the cost of originality"** — identified as the main reason for poor selection performance. The core tension is **originality vs feasibility**, and "highly original ideas tend to be disliked or rejected because they are perceived to be risky and unfeasible."
**[F]** Related: *Why Great Ideas Are Often Overlooked* (Univ. of Groningen research portal) reviews idea-evaluation/selection failure generally; Van Damme et al. (2019, *Creativity and Innovation Management*, <https://onlinelibrary.wiley.com/doi/10.1111/caim.12306>) tests epistemic/social-motivation interventions to improve group selection.

**[I] This is the most important citation in this digest for `parley-design`.** It says: *even when divergence works perfectly, the selection step regresses to the safe option.* That is the mechanism behind "AI slop" in a multi-agent setting — not that the agents can't make something distinctive, but that the **vote will kill the distinctive one**. Concrete consequences:
1. A plain majority vote over directions is the **worst** available decision rule; it is precisely the mechanism the paper indicts.
2. The Decider role exists partly to *override* feasibility bias — which is exactly why GV made the supervote **able to ignore the straw poll**.
3. Our rubric (§4) must **score originality/distinctiveness as its own weighted criterion**, or feasibility bias enters through the back door of "usability".

## 3.4 The Decider / DRI

**[F]** Apple's **DRI ("Directly Responsible Individual")** — one named person owns each significant decision or deliverable; the explicit purpose is to "install a clear owner of an outcome and **reduce decision-by-committee mindset**". The same pattern appears elsewhere as **single-threaded leader / owner / sponsor** (Amazon). Sources: <https://blog.matt-rickard.com/p/directly-responsible-individuals>, <https://www.forbes.com/sites/quora/2012/10/02/how-well-does-apples-directly-responsible-individual-dri-model-work-in-practice/>, <https://tettra.com/article/directly-responsible-individuals-guide/>
**[F]** GV's Decider is the same idea, ritualised: "**Who is the decider? She should make the final call — not the group.**"
**[I]** Note the Pixar counterweight: the Braintrust is *deliberately stripped of authority* so its candor is safe, while the **director** (the DRI) keeps the decision. **Critique power and decision power must be held by different roles.** Our protocol must not let the critics also be the deciders — otherwise critique becomes campaigning.

---

# 4. Judging rubrics — what is actually scoreable

## 4.1 Awwwards — a real weighted rubric with a real aggregation rule ★

Source: <https://www.awwwards.com/about-evaluation/>

**[F] Weights:** **Design 40 · Usability 30 · Creativity 20 · Content 10** (points, summing to 100).
**[F] Jury:** approved sites go to a **minimum of 18 jury members**.
**[F] Aggregation:** "**the 3 jury scores furthest from the average are automatically eliminated by our system**" — i.e. a **trimmed mean that drops the 3 largest absolute deviations**.
**[F] Thresholds:** **≥ 6.5 → Honorable Mention**; Site of the Day = highest scoring; **Site of the Month** = top eight of the month, subject to a **second jury review**; **Developer Award** requires SOTD **> 7.0**. Voting window: 5 days.
**[F]** The page gives weights but **no verbatim definitions** of the four criteria — the rubric is weighted but **not anchored**.

**[I]** Copy the *shape* (weighted criteria + trimmed aggregation + numeric threshold + escalation review), not the weights. And fix the gap: **add anchors.** An unanchored 0–10 from an LLM is noise (see §5.4).

## 4.2 Red Dot: Product Design — a criteria *list* (no weights)

Source: <https://www.red-dot.org/pd/faq>, <https://en.wikipedia.org/wiki/Red_Dot_Design_Award>

**[F] Nine criteria:** degree of innovation · ergonomics · product periphery · functionality · durability · **self-explanatory quality** · **formal quality** · ecological compatibility · **symbolic and emotional content**.
**[F] Jury:** ~40 international designers/professors/journalists (≈30 for Brands & Communication Design); multi-day process; **"products are not compared directly with each other but instead evaluated on an individual basis."**
**[F] Conflict rules:** jurors must be independent, cannot be employed by an industrial producer, and **cannot adjudicate designs they were involved in**.

**[I]** Two transferable rules: (a) **absolute scoring, not pairwise comparison** — this is a direct mitigation for the LLM position bias in §5.4; (b) **recusal** — an agent must not score its own direction. Both are one-line rules.

## 4.3 A' Design Award — anonymity + normalisation

Source: <https://competition.adesignaward.com/ada.php?ID=236>

**[F]** Entries are assessed "on intrinsic design quality, creativity, functionality, originality and innovation, **anonymously**". Jurors are "**compartmentalized and evaluate submissions individually in isolation**". Methods named: "**statistical score normalization**, pattern analysis, and juror verification". A free preliminary review round precedes official voting.
**[F]** No public scale, jury size, or weights — the rubric itself is not published.

## 4.4 Nielsen — the only rubric here with real anchors ★

Sources: <https://www.nngroup.com/articles/how-to-rate-the-severity-of-usability-problems/>, <https://www.cs.princeton.edu/courses/archive/spring13/cos436/handouts/Heuristics_explained.pdf>

**[F] Severity scale 0–4, verbatim anchors:**
- **0** — "I don't agree that this is a usability problem at all"
- **1** — "Cosmetic problem only: need not be fixed unless extra time is available on project"
- **2** — "Minor usability problem: fixing this should be given low priority"
- **3** — "Major usability problem: important to fix, so should be given high priority"
- **4** — "Usability catastrophe: imperative to fix this before product can be released"

**[F] Severity = f(frequency, impact, persistence):** how often it occurs; how hard it is for the user to overcome; whether it is one-time or recurring.
**[F] Evaluator count:** the cost–benefit sweet spot is **3–5 evaluators**; **5 evaluators find up to ~75% of all issues.**

**[I]** This is the model rubric for `parley-design-check`: **each level has a decision-relevant verbal anchor**, so two evaluators can converge without agreeing on taste. Blockers are only level 4 (and optionally 3). Copy the 0–4 scale *with the anchors rewritten for design-system violations*, and copy the 3–5 evaluator target — which happens to equal our roster size (claude, codex, hermes, kimi = 4). **[F]** That is a coincidence worth noting, not a derivation.

## 4.5 Calibration technique worth stealing

Source: <https://www.acceptmission.com/blog/scoring-calibration/>
**[F]** "Aligning reviewers only on the final composite score **hides where disagreement actually lives**" — calibration must be on *what criteria mean*, so "disagreements become signal rather than noise."
**[F]** **Anchor examples** are named "the most underused tool in calibration": an example labelled with its correct score *and the explanation of why it earns that score*, moving reviewers "from inference to observation."
**[I]** → Every criterion in our rubric ships with **one worked anchor example per level**, in the skill's reference files. This is the cheapest single accuracy improvement available to us.

---

# 5. LLM evidence: does multi-agent critique actually help — and for what?

## 5.1 Multi-agent debate does NOT reliably beat single-agent baselines ★

Source: ICLR 2025 Blogposts, *Multi-LLM-Agents Debate — Performance, Efficiency, and Scaling Challenges*, <https://d2jud02ci9yv69.cloudfront.net/2025-04-28-mad-159/blog/mad/>

**[F]** Verbatim: **"current MAD frameworks fail to consistently outperform simple single-agent test-time computation strategies."**
**[F]** vs CoT: "most MAD frameworks are not able to achieve consistently better performances than CoT". vs Self-Consistency: "most MAD frameworks fail to surpass SC".
**[F]** Concrete example (GPT-4o-mini, MMLU): **CoT 80.73% · SC 82.13% · MAD 74.73%** — MAD is **6 points below CoT and 7.4 below SC**.
**[F]** Benchmarks used: MMLU, MMLU-Pro, AGI-Eval, CommonSenseQA, ARC-Challenge, GSM8k, MATH, HumanEval, MBPP — **all factual/reasoning/code, none aesthetic.**
**[F]** "increasing test-time computation does not always improve accuracy"; only EOT (GSM8k, MMLU) and MAD (HumanEval) showed scalability in token spend.
**[F]** Diagnosis: "most test cases in existing benchmarks only require a single knowledge point… **MAD degrades to an inefficient resampling method**."
**[F]** Authors' hedge: "this is by no means indicating that MAD is not worth further exploration."

**[I] Direct implication for us — and it is uncomfortable:** the strongest published evidence about multi-agent debate is *negative*, and it is entirely on **verifiable** tasks. Nobody has shown MAD improves **aesthetic** outcomes. Therefore `parley-design` must **not** be justified as "more agents → better taste". Its honest justification is different and stronger: **multi-agent is a diversity generator, not an accuracy amplifier.** The win comes from producing N genuinely different directions (which a single agent demonstrably will not, §5.3) and then selecting with a *rule*, not from debating our way to good taste.

## 5.2 Deliberation degrades facts and collapses stance

Source: Wan, Wu, Luo, Li, Wang, Chen, Kan — *The Deliberative Illusion: Diagnosing Factual Attrition and Stance Homogenization in Multi-Agent LLM Deliberation*, arXiv **2606.03032v1** (3 Jun 2026), <https://arxiv.org/pdf/2606.03032>, code <https://github.com/whr000001/DelibTrace>

**[F]** Two named pathologies: **factual attrition** (agents progressively lose/misrepresent factual content across rounds) and **stance homogenization** (positions converge toward uniformity regardless of whether the converged stance is correct).
**[F]** Finding: factual accuracy **degrades measurably across deliberation rounds** while stance diversity **decreases** — "agents adopt similar positions without necessarily improving truthfulness."
**[F]** Recommendation: **monitor factual fidelity and stance diversity independently**, and evaluate deliberation systems on **factual preservation** rather than assuming discussion improves accuracy.

**[I]** Round count is a **hazard**, not a virtue. Our protocol should cap critique at **one adversarial round** by default and require an explicit reason to add a second. And "stance diversity" is a **measurable, loggable metric** — `parley-design-check` can compute it (did the four directions stay four, or did rounds 2–3 turn them into one mush?).

## 5.3 Diversity collapse in multi-agent idea generation

Source: Chen, Tong, Yang, He, Zhang, Zou, Wang, He — *Diversity Collapse in Multi-Agent LLM Systems: Structural Coupling and Collective Failure in Open-Ended Idea Generation*, arXiv **2604.18005v2** (22 Apr 2026), <https://arxiv.org/pdf/2604.18005>
**[F]** Central finding: multi-agent LLM systems **homogenize during collaborative idea generation**; **structural coupling** between agents produces convergence rather than complementary thinking. (PDF text extraction was partial — the numeric tables did not decompress; treat the numbers as unverified, the qualitative claim as verified from title+abstract-level content.)

**[F] Corroborating single-model evidence:**
- LLMs produce "dramatically less diverse creative output" than humans: **22 LLMs vs 102 humans, effect sizes 1.4–2.2** (reported in the homogenization literature; see <https://arxiv.org/html/2601.06116v3>, *The Homogenization Problem in LLMs*).
- RLHF/alignment "forces models into **mode collapse**… converges on a narrow set of safe, homogenized attractor states"; "aligned models carry **less conceptual diversity than their base counterparts**" (ibid.).
- "Instruction tuning improves **lexical** diversity while constraining **syntactic and semantic** diversity" (ibid.).
- Human/ChatGPT comparison: *Homogenizing effect of LLMs on creative diversity*, <https://www.sciencedirect.com/science/article/pii/S294988212500091X>.
- *LLMs Exhibit Significantly Lower Uncertainty in Creative Writing Than Professional Writers*, arXiv 2602.16162.

**[I]** **This is the mechanistic definition of "AI slop".** Slop is not bad craft; it is **mode collapse toward the aligned attractor**. Two protocol consequences:
1. **Isolation during divergence is not hygiene, it is the whole product.** The moment agents see each other's drafts, structural coupling starts and the four directions become one.
2. **A distinctness gate is mandatory and mechanical**: before any critique happens, verify the N directions actually differ on named axes. If two agents both produced "dark hero + purple-blue gradient + Inter", the diverge phase *failed* and must be re-run with forced-distinct assignments — not critiqued.

## 5.4 LLM-as-judge is biased — quantified

**[F] Position bias.** Shi, Ma, Liang, Diao, Ma, Vosoughi — *Judging the Judges: A Systematic Study of Position Bias in LLM-as-a-Judge*, arXiv **2406.07791** (Jun 2024, rev. Nov 2025), <https://arxiv.org/abs/2406.07791>. Scale: **15 LLM judges, ~40 solution-generating models, 2 benchmarks (MTBench, DevBench), 22 tasks, >150,000 evaluation instances.** Three metrics: **repetition stability, position consistency, preference fairness**. Findings: position bias "**is not due to random chance and varies significantly across judges and tasks**"; **weakly** influenced by prompt-component length but **strongly** affected by the **quality gap** between the compared solutions.
**[I]** → Small quality gaps (which is exactly the case when four good directions compete) are **where position bias is worst**. So: **do not use pairwise "which is better, A or B"** for the final call. Use **absolute rubric scoring** (Red Dot's model, §4.2) and randomise/rotate presentation order.

**[F] Self-preference bias.** *Self-Preference Bias in LLM-as-a-Judge*, arXiv 2410.21819, <https://arxiv.org/pdf/2410.21819>. Measured self-preference on **ArenaHard ranges −38% to +90%**; on other datasets **−21% to +56%**.
**[I]** → An agent must never score its own direction (recusal, §4.2), and the winner must not be chosen by simple self-report. Our roster is heterogeneous (claude/codex/hermes/kimi), which helps, but does not remove it.

**[F] Verbosity bias / divergence from humans.** Judge verbosity ratings track **raw response length at r = .87 vs .44 for human raters**; "judge–self-report correlations can appear to validate a construct while judge–human correlations on the same samples show no such support" (<https://arxiv.org/pdf/2606.09843>).
**[I]** → **Cap and normalise artifact length** before scoring, or the most verbose direction wins on length alone. A hard word/section cap on the direction artifact is a one-line fix with real effect.

**[F] Overconfidence.** *Trust or Escalate: LLM Judges with…* (ICLR 2025) — "existing methods for confidence estimation are **brittle even with the strongest judge model**, as they tend to **overestimate human agreement**"; simulating diverse annotators via in-context learning improves calibration and failure prediction. <https://proceedings.iclr.cc/paper_files/paper/2025/file/08dabd5345b37fffcbe335bd578b15a0-Paper-Conference.pdf>

## 5.5 The ceiling: humans barely agree with each other on aesthetics ★

Source: Zhang, Li, Wang, Shen, Hu — *Beauty in the Eye of AI: Aligning LLMs and Vision Models with Human Aesthetics in Network Visualization*, arXiv **2604.03417v1**, <https://arxiv.org/html/2604.03417>

**[F]** Human–human pairwise alignment on aesthetic choices — "the fraction of graphs on which two labelers make the same choice", micro-averaged over all overlapping labeler pairs — is **38.34%**.
**[F]** GPT-4o-mini + DINOv2 memory bank: **33.11%**. Fine-tuned DINOv2 vision model: **36.81%**. With a confidence threshold retaining ~65% of labels, the LLM reached **38.34%** — i.e. **human-level, but only on the confident subset**.
**[F]** Authors' conclusion: "AI can perform at a level that makes it an **effective proxy for human labelers**."

**[I] The single most important number in this digest.** Human inter-rater agreement on aesthetics is **~38%** — barely above the chance floor of the task. Consequences that must be written into `parley-design-check`:
1. **An aesthetic score is not a gate.** You cannot block a build on a judgement whose human ceiling is 38% agreement. Gates must be on **mechanical, checkable** properties (contrast ratios, token usage, scale conformance, banned-value lists) — Hallmark's 58-gate model — while taste stays **advisory**.
2. **Confidence-thresholded abstention is legitimate and evidenced:** the judge reached human parity *only* by abstaining on ~35% of cases. Our judge should be allowed to return `ABSTAIN` and escalate to the human Decider. That is not a cop-out; it is the published mechanism for reaching human-level.
3. Any claim in our skill that "agents will converge on good taste" is **contradicted by evidence** and must not be written.

---

# 6. Reconciling our proposed method against the evidence

| Our move | Human precedent | Verdict |
|---|---|---|
| **DIVERGE — each agent a different direction, in isolation** | GV four-step sketch "work alone together"; NN/g charrette "no talking once sketching begins"; Diehl & Stroebe nominal>interacting (r≈.56–.57) | **Strongly grounded.** Reinforced by LLM-native diversity-collapse evidence (§5.3). Strengthen: add a mechanical distinctness gate. |
| **Adversarial cross-critique** | GV Speed Critique (3 min/sketch, author silent); Figma Silent Critique; Braintrust; Black hat by assignment | **Grounded in form.** But cap rounds — deliberation causes factual attrition + stance homogenization (§5.2). Critique ≠ debate to consensus. |
| **ONE direction wins whole** | GV Supervote (Decider, 3 votes, may ignore the poll); Apple DRI; Braintrust director owns the fix | **Grounded** — provided the decider is a *role*, not a vote, and preferably the **human**. Majority vote is the failure mode (§3.3). |
| **Graft 2–3 details from losers** | GV **All-in-One** ("put them together into your best shot"); Heat Map dot clusters name the graftable parts; d.school "I like"; "maybe-laters" are kept | **Grounded**, but currently under-specified: we have no mechanism that *identifies which details are graftable*. Heat Map does exactly that. Adopt it. |
| **(missing) Rumble** | GV: conflicting winners → build both, decide by external evidence + fake brands | **Gap.** We have no branch for "two directions are incommensurable". Add it (§7 P6b). |
| **(missing) Intra-agent divergence** | Crazy 8s: 8 variations of *your own* best idea, 1 min each | **Gap.** Cheap to add, widens each agent's own search before lock-in. |

---

# 7. `DCP/1` — the Divergent Critique Protocol (concrete, ritualised, protocol-shaped)

**[I] Everything in §7 is my synthesis.** Each phase cites the human precedent it derives from. Written protocol-style (typed artifacts, phase IDs, MUST/SHOULD, versioned) per the owner's AG-UI preference. "Budget" replaces wall-clock timeboxes because agents have no clock — but the *ratios* are preserved from the human originals.

**Version:** `DCP/1` · **Roles:** `Proposer` (each roster agent) · `Critic` (each roster agent, on others' work only) · `Facilitator` (the driver/orchestrator; deterministic, non-creative) · `Scribe` (the orchestrator; writes the ledger) · **`Decider` (exactly one; the human by default)**.

**Invariant D-0 (authority split).** *[from Pixar Braintrust §1.3 + GV Note-and-Vote step 6 + Apple DRI §3.4]* — Critics MUST have **no** decision authority. The Decider MUST NOT participate as a Critic. Votes are **advisory only** and MUST be labelled as such in the ledger.

**Invariant D-1 (recusal).** *[Red Dot §4.2, self-preference bias §5.4]* — No agent may score, rank, or vote for its own direction. Its self-scores MUST be discarded, not down-weighted.

**Invariant D-2 (absolute, not pairwise).** *[Red Dot "not compared directly with each other"; position bias §5.4 worst at small quality gaps]* — Scoring MUST be absolute against the rubric. Pairwise "A or B?" prompts are FORBIDDEN in the deciding phase.

**Invariant D-3 (length normalisation).** *[verbosity bias r=.87, §5.4]* — Direction artifacts MUST obey a hard section/word cap. Over-cap artifacts are truncated before scoring, not rewarded.

---

### **P0 · BRIEF** — artifact `BRIEF.md`
*[GV "set the stage" §1.1; Irizarry "set the foundation" §1.4; Google's single signup doc §1.5]*
Facilitator writes, Decider ratifies. MUST contain: the problem; the audience; **business + user goals critique will be judged against**; hard constraints (brand, a11y target, tech); **fidelity level** expected; and the **`divergence_axes`** — the named axes on which directions MUST differ (e.g. `structure`, `typographic voice`, `colour strategy`, `motion posture`, `density`).
Gate **G0**: no `divergence_axes` → no P1.

### **P1 · DIVERGE (isolated)** — artifact `directions/<agent>/DIRECTION.md`
*[GV four-step sketch §2.2; NN/g charrette §2.1; Diehl & Stroebe §3.2; diversity collapse §5.3]*
Each Proposer works **without reading any other agent's output**. Isolation is a MUST, mechanically enforced (separate dirs; no cross-reads until P2 opens).
Sub-steps, mirroring the four-step sketch:
- **P1.a Notes** — harvest from `BRIEF.md` + any provided references. Budget 20%.
- **P1.b Ideas** — private candidate list; circle the most promising. Budget 20%.
- **P1.c Crazy-8** — **8 one-line variations of your own best idea**, no elaboration. Budget 10%. *[Crazy 8s §2.2]*
- **P1.d Direction Sketch** — the committed artifact. Budget 50%.
`DIRECTION.md` MUST be **self-explanatory without its author present** *[GV solution-sketch rule]*, MUST carry a **catchy one-word handle** (not the agent's name), and MUST declare its position on every `divergence_axis`. It MUST include a token table (type scale, colour roles, spacing scale, radius, motion durations) so it is checkable, not just describable.

### **P2 · EXHIBIT (anonymise + distinctness gate)** — artifact `EXHIBIT.md`
*[GV Art Museum §2.3 — one row, no presentation, anonymous]*
Facilitator (deterministic, no model call) strips author identity, assigns stable random slugs (`A`,`B`,`C`,`D`), **randomises order per reviewer** *[position bias §5.4]*, and publishes all directions at once.
Gate **G1 — DISTINCTNESS** *[§5.3]*: if any two directions match on **all** `divergence_axes`, or share a banned-slop signature (identical font stack + identical colour strategy + identical macrostructure), P1 **FAILED**. Re-run P1 with **forced-distinct axis assignments** (each agent is *assigned* a different position on the primary axis). MUST NOT proceed to critique on a collapsed set — critiquing a collapsed set launders the collapse into a "consensus".

### **P3 · HEAT MAP (silent, parallel, part-level)** — artifact `HEATMAP.jsonl`
*[GV Heat Map §2.3; Figma Silent Critique §1.2]*
Every Critic, **independently and without seeing others' marks**, emits typed records against *parts* of *other agents'* directions:
```
{ "direction": "C", "part": "nav.sticky-condense", "mark": "like",
  "intensity": 2, "note": "..." }        // intensity 1-3, budget 20-30 marks total
{ "direction": "C", "part": "hero.gradient", "mark": "concern",
  "severity": 3, "note": "..." }         // severity uses the 0-4 anchors, §4.4
```
`mark: "like"` records are the **graft harvest** — this phase exists as much to find graftable parts as to find flaws. *[d.school "I like" §1.6; GV "maybe-laters" §2.3]*
Reveal is **simultaneous**: no Critic sees another's marks until all are submitted. *[Note-and-Vote "commit your vote to paper" §2.4]*

### **P4 · SPEED CRITIQUE (one round, assigned lenses)** — artifact `CRITIQUE.md`
*[GV Speed Critique §2.3; Six Thinking Hats §1.7; Braintrust §1.3; Irizarry §1.4]*
- Facilitator narrates each direction from the artifact and names the **dot clusters** first *[GV step ②→③]*.
- Each Critic is **assigned a lens** (Black = risk/failure, Yellow = value/upside, White = facts/constraints/a11y, Green = adjacent alternatives) — assigned, not chosen, because the artificiality is what separates ego from the critique *[De Bono, verbatim §1.7]*.
- **The author is silent during critique of its own direction** and gets a single **`REBUTTAL.md`** afterwards, addressing only misreadings *[GV step ③.5–6]*.
- Every critique entry MUST be typed:
  `{ target, part, class: like|wish|what_if, severity: 0-4, tied_to_goal: <BRIEF goal id>, fix?: <string> }`
  - `class` from d.school *[§1.6]*; `severity` anchors from Nielsen *[§4.4]*.
  - **`tied_to_goal` is REQUIRED for `wish`** — "Good feedback is about how the design is meeting (or missing) the customer and business goals" *[GV, verbatim §1.1]*. A `wish` with no goal link is dropped as taste.
  - **`fix` is OPTIONAL and non-binding**, and MUST be rendered as a suggestion. This is the §1.3 reconciliation: diagnosis is owed, prescription is not owned *[Braintrust "notes not prescriptions" + GV "suggestions, not mandates"]*.
- **Hard cap: ONE round by default.** A second round requires an explicit Decider instruction and MUST log the reason. *[factual attrition + stance homogenization §5.2]*
- Facilitator MUST log **stance-diversity before vs after** the round. If diversity dropped below threshold, flag `HOMOGENIZATION_WARNING` in the ledger. *[§5.2 "monitor factual fidelity and stance diversity independently"]*

### **P5 · SCORE + STRAW POLL (advisory)** — artifact `SCORECARD.md`
*[Awwwards weighted rubric + trimmed mean §4.1; GV Straw Poll §2.3; Rietzschel §3.3]*
Each Critic scores each **other** direction absolutely (D-1, D-2) on an anchored, weighted rubric. Suggested starting weights, explicitly derived from Awwwards but **re-weighted to counteract feasibility bias** *[§3.3]*:

| Criterion | Weight | Anchored at 0/3/5/8/10 by worked examples |
|---|---|---|
| **Distinctiveness** (does it avoid the aligned attractor / slop signature?) | 30 | Y |
| **Systemic coherence** (tokens actually generate the UI; no orphan values) | 25 | Y |
| **Fitness to brief** (goals in `BRIEF.md`, explicitly) | 25 | Y |
| **Craft & accessibility** (mechanical: contrast, hit targets, focus, scale conformance) | 20 | Y |

Aggregation: **drop the highest and lowest score per direction, then mean** — the Awwwards trim, scaled to a 4-agent jury *[§4.1]*.
Then a **Straw Poll**: each Critic commits **one** vote *[GV ④]* — for a whole direction **or a single part inside one** — with a one-paragraph reason, all committed before reveal. **Explicitly advisory.**
An `ABSTAIN` verdict is legitimate and MUST be preserved, not coerced into a vote *[confidence-thresholded abstention reached human parity, §5.5]*.

### **P6 · SUPERVOTE (the decision)** — artifact `DECISION.md`
*[GV Supervote §2.3; Note-and-Vote step 6 §2.4; DRI §3.4]*
The **Decider** (human by default) holds **three supervotes** and **MAY ignore the straw poll and the scorecard entirely**. `DECISION.md` MUST record: winning direction, the losers marked **`maybe-later` (retained, not deleted)** *[§2.3]*, and — if the Decider overrode the poll — a one-line reason. Non-winning directions are **archived, not discarded**.
**No averaging.** Combining two directions' *visual systems* into a compromise is FORBIDDEN. *[§3.1 rhetoric; §3.3 evidence — averaging is exactly the feasibility-biased regression]*

### **P6b · RUMBLE (branch)** — artifact `RUMBLE.md`
*[GV Rumble/All-in-One §2.3 — the branch our method was missing]*
If the Decider judges the top two directions **genuinely incommensurable** (they answer the brief with conflicting premises, not conflicting details), the Decider MAY declare a Rumble instead of picking: **build both to a cheap, comparable fidelity**, give each a **distinct fake handle** so they are not read as "version A / version B" *[verbatim: "you risk sounding like an optometrist"]*, and defer the decision to **external evidence** (user test, stakeholder, metric).
Default is **All-in-One**: *"If you think you can combine your winning sketches into one product, don't bother with a rumble."*
Rumble MUST be rare and MUST be justified in writing — it doubles cost.

### **P7 · GRAFT (bounded, from the heat map)** — artifact `GRAFTS.md`
*[GV All-in-One §2.3; Heat Map dot clusters §2.3; d.school "I like" §1.6]*
The **winning direction's author** is the DRI for the graft *[Braintrust: the director owns the fix, §1.3]*.
- Graft candidates MUST come from `HEATMAP.jsonl` `mark: "like"` clusters on losing directions — **not** from fresh invention in this phase.
- **Maximum 3 grafts.** Each graft MUST be a **discrete, nameable detail** (an interaction, a component treatment, a copy device, a motion rule) — **never** a token-system layer (never a colour ramp, never a type scale, never a grid). Grafting a system layer is how you get a camel.
- Each graft MUST state: source direction, the exact part, why it survives inside the winner's system, and **which winner token it must be re-expressed in**. A graft that cannot be re-expressed in the winner's tokens is **rejected**.
- Gate **G2 — COHERENCE**: after grafting, re-run the mechanical checks (`parley-design-check`). Any new orphan token, off-scale value, or contrast failure fails the graft, not the winner.

### **P8 · RATIFY** — artifact `DESIGN.md` (the design system) + `LEDGER.md`
*[Hallmark's `design.md` single-source-of-truth; Google's "never established" §1.5]*
`DESIGN.md` becomes the canonical, machine-readable system; from here on, **the system is the authority, not the critique**. `LEDGER.md` records the whole run: every phase artifact hash, the scorecard, the straw poll, the supervote, the override reason, the grafts, and the `HOMOGENIZATION_WARNING` flags. The retro *[protocol §13]* MUST be able to ask: *did the direction that won actually ship, and did the grafts survive?*

---

# 8. Corrections to our stated method (the honest part)

1. **[I] "One winner whole, always" is stricter than GV.** GV branches: conflicting → **Rumble**; compatible → **All-in-One**. Adopt the branch (P6b) or we will force a premature pick on genuinely incommensurable directions. *[§2.3]*
2. **[I] "Adversarial cross-critique" must be capped at one round.** More rounds = factual attrition + stance homogenization, evidenced. Our instinct ("more debate = better") is contradicted. *[§5.2]*
3. **[I] The winner must not be chosen by agent majority vote.** Idea *selection* is where the documented failure lives — feasibility bias kills the original option. Vote is advisory; a single Decider (preferably the human) decides. *[§3.3, §2.4]*
4. **[I] We have no mechanism for choosing grafts.** Add the Heat Map — its `like` clusters *are* the graft shortlist, produced before anyone argues. *[§2.3]*
5. **[I] We are missing the distinctness gate.** Without G1, the four "different directions" will quietly be one direction, and every downstream ritual becomes theatre. *[§5.3]*
6. **[I] Do not claim multi-agent improves taste.** The published MAD evidence is negative and entirely on verifiable tasks; human aesthetic agreement is ~38%. Claim what is defensible: multi-agent **generates diversity** that a single aligned model will not, and a **rule** (not a debate) selects from it. *[§5.1, §5.5]*

---

# 9. FACTS vs INFERENCE — index

**Verified facts** (all cited inline): GV Sticky Decision 5 sub-steps incl. 20–30 dots, 3 min/sketch, author silent, 1 straw vote, 3 supervotes, Decider may ignore the poll · Rumble vs All-in-One rule and fake brands · Note-and-Vote 7 steps with 5–10/2/5 min boxes and "no sales pitch" · four-step sketch 20/20/8/30–90 min, self-explanatory + anonymous, "individuals working alone generate better solutions than groups brainstorming out loud" · NN/g charrette 5 min sketch / 2 min show / 1 min Q&A, "no talking once sketching begins" · Figma's six formats, RTR 2 min/person, max two topics per hour · GV critique five steps + five verbatim rules · Braintrust "no authority" + "notes not prescriptions" · plussing "criticism must come with a plus" · Irizarry's 12 rules · Six Thinking Hats roles + "blatantly artificial… separate individual ego" · Awwwards 40/30/20/10, ≥18 jurors, drop 3 furthest from average, ≥6.5 HM, >7.0 Developer Award · Red Dot 9 criteria, ~40 jurors, "not compared directly with each other", recusal · A' anonymity + compartmentalised jurors + score normalization · Nielsen 0–4 anchors, frequency/impact/persistence, 3–5 evaluators ≈75% · Diehl & Stroebe 1987 production blocking, r=.57/.56 · Rietzschel et al. 2010 feasibility-over-originality · camel proverb → Reader's Digest Sept 1954, T. R. Quaife, Kettering/Issigonis unsupported · MAD: CoT 80.73 / SC 82.13 / MAD 74.73 on MMLU, "fail to consistently outperform" · arXiv 2606.03032 factual attrition + stance homogenization · arXiv 2604.18005 diversity collapse · arXiv 2406.07791 15 judges / 150k instances / bias strongest at small quality gaps · self-preference −38%..+90% ArenaHard · verbosity r=.87 vs .44 · aesthetics human–human agreement **38.34%**, LLM 33.11%, VM 36.81%, LLM reaches 38.34% with ~65% coverage.

**Inference (mine, unverified):** the Braintrust/plussing contradiction and its phase-based reconciliation · Silent Critique as the right analogue for headless agents · "slop = mode collapse toward the aligned attractor" · the whole of §7 (`DCP/1`) · the six corrections in §8 · the re-weighted rubric in P5 · "grafts must never be system layers" · treating ~38% agreement as a hard ceiling that forbids aesthetic gating.

**Could not verify:** the "three concepts" agency convention (no authoritative source — folklore) · the *Sprint* book's exact page-level wording (only secondary summaries reachable) · Catmull's verbatim *Creativity, Inc.* text (Fast Company excerpt returned HTTP 403) · A' Design Award's actual scale/weights (unpublished) · numeric tables in arXiv 2604.18005 (PDF streams did not extract).

---

## Transferable to `parley-design` / `parley-design-check`

Ranked by value × cost-to-implement.

1. **The Sticky Decision, transposed whole (Art Museum → Heat Map → Speed Critique → Straw Poll → Supervote).** *[§2.3]* This is a battle-tested, named, five-step ritual that matches our intended flow almost move-for-move, and it comes with a decision rule that is explicitly *not* a vote. Adopt the names verbatim — named rituals are memorable and auditable, and the owner wants protocol-shaped specs. → `parley-design` phases P2–P6.
2. **G1 — the DISTINCTNESS GATE before critique.** *[§5.3]* Highest-value *new* mechanism. Multi-agent LLM systems homogenize; if the four directions collapse, everything downstream is theatre. Mechanical, cheap, and it is the one gate that protects the entire premise of the skill. → `parley-design-check` ships the checker; `parley-design` ships the axes.
3. **Authority split: Critics have no authority; exactly one Decider; votes are advisory.** *[Braintrust §1.3 + GV Decider §2.4 + DRI §3.4 + Rietzschel §3.3]* One paragraph of doctrine that prevents the single most documented failure (selection regressing to the safe option).
4. **Nielsen's 0–4 severity scale with verbal anchors, re-anchored for design-system violations, + "3–5 evaluators find ~75%".** *[§4.4]* The only rubric in this research with real anchors. It is what makes `parley-design-check` mechanical rather than opinionated, and it justifies our roster size.
5. **Typed critique records: `{class: like|wish|what_if, severity, tied_to_goal, fix?}`.** *[d.school §1.6 + GV §1.1 + Braintrust/plussing §1.3]* Converts prose critique into machine-checkable data. Two specific rules: a `wish` **must** cite a brief goal or it is dropped as taste; a `fix` is **optional and non-binding**.
6. **Heat Map `like` clusters as the graft shortlist, + max 3 grafts, + never graft a system layer.** *[§2.3, §7 P7]* Gives the graft step a mechanism it currently lacks, and the "never a token layer" rule is the concrete anti-camel constraint.
7. **Anonymity + randomised order + recusal + absolute (not pairwise) scoring + length caps.** *[§2.2, §4.2, §5.4]* Four one-line rules that each neutralise a *quantified* LLM judging bias. Cheapest evidence-to-effort ratio in the whole digest.
8. **Awwwards' aggregation shape: weighted criteria + trimmed mean (drop extremes) + numeric threshold + escalation to a second review.** *[§4.1]* Copy the shape, replace the weights (bump Distinctiveness to counteract feasibility bias), and add the anchors Awwwards lacks.
9. **ONE critique round by default; log stance-diversity before/after; `HOMOGENIZATION_WARNING`.** *[§5.2]* Counter-intuitive and therefore worth writing down explicitly, plus it is a metric `parley-design-check` can actually compute.
10. **The Rumble branch.** *[§2.3]* The missing escape hatch for incommensurable directions: build both cheap, fake-brand them, decide by external evidence. Default remains All-in-One.
11. **`ABSTAIN` is a legitimate verdict; aesthetic judgements are advisory, never gates.** *[§5.5 — 38.34% human ceiling; confidence thresholding is how the judge reached parity]* This is the intellectual honesty that will make the skill credible: gate on mechanics, advise on taste.
12. **Crazy-8s as an intra-agent divergence step** (8 one-line variations of your own best idea before committing). *[§2.2]* Cheap, timeboxed, widens each agent's own search.
13. **Assigned critique lenses (Black/Yellow/White/Green), justified by De Bono's "blatantly artificial… separates ego from the activity".** *[§1.7]* Gives us a principled reason to *assign* adversarial roles instead of asking for honest opinions.
14. **`BRIEF.md` must name the `divergence_axes` before P1.** *[§1.1 "set the stage" + §1.5 preparation questions]* Without declared axes, "be different" is unenforceable and G1 is uncheckable.
15. **"Never established"** *[§1.5]* — a standing tenet that the protocol's own format is subject to revision. Pairs naturally with our §13 retro.
16. **Provenance discipline for the camel proverb** *[§3.1]* — if we quote it, quote Quote Investigator's provenance too. A skill that spreads a misattribution while preaching rigour is self-refuting.

## Do NOT copy

1. **"Three concepts" as a rule.** No authoritative source exists — it is agency folklore (§2.5). Deriving our N from it would put an unsourced claim at the foundation of the skill. Derive N from roster size + Nielsen's 3–5 evaluators instead.
2. **The camel proverb as evidence.** It is a 1954 magazine quip with a contested attribution (§3.1). Use Rietzschel et al. 2010 (§3.3) as the actual evidence for the same point.
3. **Diehl & Stroebe's *mechanism*.** Production blocking is "can't generate while listening" — a serial-speech constraint headless agents do not have (§3.2). Borrow the conclusion, cite the LLM-native mechanism (diversity collapse) for the *reason*.
4. **Spoken-crit turn-taking machinery** — Round-The-Room, Popcorn, 2-min-per-person, "max two topics per hour" (§1.2). These solve one-mouth-at-a-time and meeting fatigue. Agents have neither. Copying them adds ceremony with zero effect. Use Figma's **Silent Critique** instead.
5. **Wall-clock timeboxes (5 min, 3 min, 8 min).** Preserve the *ratios* (they encode how much effort each step deserves), discard the minutes. Agents have no clock; a "3-minute critique" is meaningless and will be ignored or hallucinated as compliance.
6. **Plussing as a blanket rule ("every criticism must come with a fix").** It contradicts notes-not-prescriptions (§1.3) and, worse, in an LLM setting it pushes critics to *author* the fix — which is how a critique quietly becomes an unratified redesign and how stance homogenization starts. Keep `fix` optional and explicitly non-binding.
7. **Pairwise "which is better, A or B?" for the final call.** Position bias is real, non-random, and **strongest exactly when the quality gap is small** (§5.4) — which is our situation with four good directions. Absolute rubric scoring only.
8. **Letting an agent score its own direction, even down-weighted.** Self-preference ranges −38% to +90% (§5.4). Discard, do not weight.
9. **Multi-round debate-to-consensus.** Evidenced to lose facts and collapse stances (§5.2), and MAD underperforms plain self-consistency on the tasks where it has been measured (§5.1). One round, then decide.
10. **Any aesthetic score used as a hard gate.** Human–human agreement on aesthetics is 38.34% (§5.5). A blocking gate built on that is a coin flip wearing a lab coat. Gate on mechanics; advise on taste.
11. **"More agents → better taste" as the skill's justification.** Directly contradicted by §5.1. Justify on diversity generation + rule-based selection instead — a claim that survives the evidence.
12. **Averaging or blending two directions' token systems into a compromise.** No source supports it; §3.3 explains why it regresses to safe. Grafts are discrete details only, re-expressed in the winner's tokens.
13. **Red Dot / A' Design Award as rubric sources.** Red Dot publishes criteria but no weights, no scale, no anchors; A' publishes neither. Borrow their *procedural* rules (absolute scoring, recusal, anonymity, normalisation) — not their rubrics, which are not actually specified.
14. **Awwwards' weights as-is (Design 40 / Usability 30 / Creativity 20 / Content 10).** They are tuned for judging finished award-bait websites, and they under-weight the thing we most need to protect (distinctiveness at 20). Copy the *aggregation*, re-derive the weights.
