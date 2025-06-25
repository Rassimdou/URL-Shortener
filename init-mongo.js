db = db.getSiblingDB('urlshortener');

// Create collections
db.createCollection('urls');
db.createCollection('stats');

// Create indexes
db.urls.createIndex({ shortCode: 1 }, { unique: true });
db.urls.createIndex({ expiresAt: 1 }, { expireAfterSeconds: 0 });
db.stats.createIndex({ _id: 1 });

// Insert initial stats record if it doesn't exist
if (db.stats.countDocuments({ _id: "totalClicks" }) === 0) {
    db.stats.insertOne({ _id: "totalClicks", value: 0 });
}