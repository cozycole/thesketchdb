import { db } from "../db";

export const seedSketches = () => {
  if (db.sketches.count() > 0) return;

  const sketchTitles = [
    "Office Icebreaker Goes Off the Rails",
    "Zoom Meeting Nobody Can Leave",
    "HR Training Video From Hell",
    "Performance Review With Too Much Honesty",
    "Coworker Who Treats Slack Like Twitter",
    "Company Retreat That Turns Into a Cult",
    "Team Building Exercise Gets Too Real",
    "Manager Who Just Discovered Buzzwords",
    "Intern Who Thinks They’re the CEO",
    "Office Birthday Party Nobody Asked For",

    "Awkward First Date at a Coffee Shop",
    "Dating App Bio Consultant",
    "Couple Arguing Over Split Checks",
    "Roommate Interview That Feels Like a Job",
    "Breakup That Turns Into a PowerPoint",
    "Friends Who Overshare in Public",
    "Double Date With Unequal Couples",
    "Wedding Toast That Never Ends",
    "Engagement Announcement Nobody Supports",
    "Couple That Calls Each Other ‘Best Friend’",

    "Family Dinner After Watching the News",
    "Parents Trying to Be Relatable",
    "Sibling Rivalry at Thanksgiving",
    "Dad Who Just Learned What a Podcast Is",
    "Mom Who Treats Google Like a Person",
    "Family Group Chat Meltdown",
    "Holiday Gift Exchange Gone Wrong",
    "Cousin With a New Pyramid Scheme",
    "Grandparent Who Believes Everything Online",
    "Family Game Night Turns Competitive",

    "Local News Anchor Loses Control",
    "Morning Show Cooking Segment Disaster",
    "True Crime Podcast About Something Mundane",
    "Documentary Narrator Overdramatizes Everything",
    "Reality Show Contestants Immediately Cry",
    "Streaming Service Original Trailer",
    "Movie Trailer That Reveals the Entire Plot",
    "Award Show Acceptance Speech Apocalypse",
    "Late Night Monologue Writer Strike",
    "Celebrity Apology Tour",

    "Political Ad Nobody Understands",
    "Town Hall Meeting Goes Sideways",
    "Campaign Staffer Who Can’t Stop Tweeting",
    "Debate Moderator Loses Authority",
    "Focus Group With Wild Opinions",
    "Press Secretary Dodges Every Question",
    "Grassroots Movement With No Plan",
    "Local Election Nobody Prepared For",
    "Public Forum Becomes Group Therapy",
    "Voter Guide That Explains Nothing",

    "Tech Startup Pitch That Makes No Sense",
    "App That Does Too Much",
    "Software Update Ruins Someone’s Life",
    "Tech Support Call From Another Dimension",
    "AI Assistant With Too Much Personality",
    "Password Requirements Gone Insane",
    "Beta Feature Accidentally Goes Live",
    "Product Launch With Immediate Apology",
    "Social Media Platform Rebrand Disaster",
    "Influencer Explaining Blockchain",

    "Fitness Class That Feels Like Boot Camp",
    "Yoga Instructor Overshares",
    "Personal Trainer Who Judges Silently",
    "Meditation App That Causes Stress",
    "Wellness Retreat Turns Competitive",
    "Diet Trend Nobody Can Explain",
    "Gym Equipment Nobody Knows How to Use",
    "Spin Class With Aggressive Motivation",
    "Health Podcast Host Isn’t Healthy",
    "Self-Care Routine Takes All Day",
  ];

  let idCount = 1;
  sketchTitles.forEach((title, i) => {
    db.sketches.create({
      id: idCount++,
      title: title,
      slug: title.toLowerCase().replace(/[^a-z0-9]+/g, "-"),
      url: `https://example.com/sketch-${i + 1}`,
      uploadDate: new Date(2018 + (i % 6), i % 12, (i % 28) + 1),
      popularity: Math.floor(Math.random() * 1000),
      rating: Number((Math.random() * 5).toFixed(1)),
    });
  });
};
