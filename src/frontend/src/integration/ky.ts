import ky from "ky";

const kyClient = ky.create({
  prefixUrl: "/api",
  timeout: 30000,
  retry: {
    limit: 2,
    methods: ["get"],
    statusCodes: [408, 413, 429, 500, 502, 503, 504],
  },
});

export default kyClient;
