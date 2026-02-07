import { db } from "../db";

export const seedShows = async () => {
  if (db.shows.count() > 0) return;

  const shows = [
    {
      name: "Saturday Night Live",
      profileImage:
        "https://thesketchdb-testing.nyc3.cdn.digitaloceanspaces.com/show/small/snl.jpg",
      url: "https://www.youtube.com/SaturdayNightLive",
      alias: "snl",
    },
    {
      name: "Whitest Kids U' Know",
      profileImage:
        "https://thesketchdb-testing.nyc3.cdn.digitaloceanspaces.com/show/small/wkuk.jpg",
      url: "https://www.youtube.com/@OfficialWKUK",
      alias: "wkuk",
    },
  ];

  let showIdCount = 1;
  let seasonIdCount = 1;
  let episodeIdCount = 1;

  for (const [i, show] of shows.entries()) {
    const showRecord = await db.shows.create({
      id: showIdCount++,
      slug: show.name.toLowerCase().replace(/[^a-z0-9]+/g, "-"),
      name: show.name,
      url: show.url,
      alias: show.alias,
      profileImage: show.profileImage,
      establishedDate: new Date(2018 + (i % 6), i % 12, (i % 28) + 1),
    });

    const seasonCount = 3;
    const epCount = 5;

    for (let j = 0; j < seasonCount; j++) {
      const seasonRecord = await db.seasons.create({
        id: seasonIdCount++,
        slug: String(showRecord.slug) + `-s${j + 1}`,
        seasonNumber: j + 1,
        show: showRecord,
      });

      for (let k = 0; k < epCount; k++) {
        await db.episodes.create({
          id: episodeIdCount++,
          slug: String(showRecord.slug) + `-s${j + 1}-e${k + 1}`,
          season: seasonRecord,
          episodeNumber: k + 1,
        });
      }
    }
  }
};
