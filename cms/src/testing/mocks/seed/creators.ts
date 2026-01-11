import { db } from "../db";

export const seedCreators = async () => {
  if (db.creators.count() > 0) return;

  const creators = [
    {
      name: "Almost Friday TV",
      profileImage:
        "https://thesketchdb-testing.nyc3.cdn.digitaloceanspaces.com/creator/small/almostfriday.jpg",
      url: "https://www.youtube.com/@AlmostFridayTV",
      alias: "",
    },
    {
      name: "Gilly and Keeves",
      profileImage:
        "https://thesketchdb-testing.nyc3.cdn.digitaloceanspaces.com/creator/small/gilly-and-keeves.jpg",
      url: "https://www.youtube.com/c/gillyandkeeves",
      alias: "gilly",
    },
    {
      name: "College Humor",
      profileImage:
        "https://thesketchdb-testing.nyc3.cdn.digitaloceanspaces.com/creator/small/collegehumor.jpg",
      url: "https://www.youtube.com/collegehumor",
      alias: "dropout",
    },
  ];

  let idCount = 1;
  creators.forEach((creator, i) => {
    db.creators.create({
      id: idCount++,
      slug: creator.name.toLowerCase().replace(/[^a-z0-9]+/g, "-"),
      name: creator.name,
      url: creator.url,
      alias: creator.alias,
      profileImage: creator.profileImage,
      establishedDate: new Date(2018 + (i % 6), i % 12, (i % 28) + 1),
    });
  });
};
