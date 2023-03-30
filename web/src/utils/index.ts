// regular expression to validate permalink with
const permalinkRegex =
  /https:\/\/github\.com\/[a-zA-Z0-9-_\.]+\/[a-zA-Z0-9-_\.]+\/blob\/[a-z0-9]{40}(\/[a-zA-Z0-9-_\.]+)+#L[0-9]+-L[0-9]+/;

// function to get the permalink
export const getPermalinkPreview = async (
  permalink: string
): Promise<string|null> => {
  try {
    // test the input returning null if not valid
    if (!permalinkRegex.test(permalink)) {
      return null;
    }

    // parse the permalink
    const permalinkUrl = new URL(permalink);

    // get the range start/end from the hash
    const [rangeStart, rangeEnd] = permalinkUrl.hash
      .slice(1) // remove the '#'
      .split("-") // separate start/end
      .map((v) => parseInt(v.slice(1))); // remove the 'L' from start/end & parse as int
    permalinkUrl.hash = "";

    // change the host from github.com to raw.githubusercontent.com
    permalinkUrl.host = "raw.githubusercontent.com";

    // remove the /blob segment from the url
    permalinkUrl.pathname = permalinkUrl.pathname
      .split("/")
      .filter((part) => part != "blob")
      .join("/");

    const response = await fetch(permalinkUrl);
    const contents = await response.text();
    const contentRange = contents
      .split(/\r\n|\n|\r/)
      .slice(rangeStart - 1, rangeEnd)
      .join("\n");
    return contentRange;
  } catch (e: any) {
    //TODO: better error handling
    console.error("Failed to fetch the code preview", e);
    return null;
  }
};
