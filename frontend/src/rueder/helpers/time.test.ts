import * as time from "./time"

const nowStr = "2022-02-22 22:02:22"
const now = new Date(nowStr)
const todayStr = "2022-02-22 10:10:01"
const yesterdayStr = "2022-02-21 10:10:01"
const lastYearStr = "2021-02-22 22:02:22"

describe("localizeTime", () => {
    it("builds time string for 'now'", () => {
        const result = time.localizeTime(nowStr, false, now)
        expect(result.isToday).toBeTruthy()
        expect(result.date).toBe("22.02.2022")
        expect(result.time).toBe("22:02")
        expect(result.toString()).toBe("22.02.2022 22:02")
    })

    it("builds time string for another time today", () => {
        const result = time.localizeTime(todayStr, false, now)
        expect(result.isToday).toBeTruthy()
        expect(result.date).toBe("22.02.2022")
        expect(result.time).toBe("10:10")
        expect(result.toString()).toBe("22.02.2022 10:10")
    })

    it("builds time string with hideYear", () => {
        const result = time.localizeTime(todayStr, true, now)
        expect(result.isToday).toBeTruthy()
        expect(result.date).toBe("22.02.")
        expect(result.time).toBe("10:10")
    })

    it("builds time string with 'yesterday'", () => {
        const result = time.localizeTime(yesterdayStr, true, now)
        expect(result.isToday).toBeFalsy()
        expect(result.date).toBe("21.02.")
        expect(result.time).toBe("10:10")
    })

    it("builds time string with 'last year' and hideYear=true", () => {
        const result = time.localizeTime(lastYearStr, true, now)
        expect(result.isToday).toBeFalsy()
        expect(result.date).toBe("22.02.2021")
        expect(result.time).toBe("22:02")
        expect(result.toString()).toBe("22.02.2021 22:02")
    })

    it("builds time string with 'last year' and hideYear=false", () => {
        const result = time.localizeTime(lastYearStr, false, now)
        expect(result.isToday).toBeFalsy()
        expect(result.date).toBe("22.02.2021")
        expect(result.time).toBe("22:02")
        expect(result.toString()).toBe("22.02.2021 22:02")
    })

    it("builds time string with 'last year' and hideYear=false", () => {
        const result = time.localizeTime(lastYearStr, false, now)
        expect(result.isToday).toBeFalsy()
        expect(result.date).toBe("22.02.2021")
        expect(result.time).toBe("22:02")
        expect(result.toString()).toBe("22.02.2021 22:02")
    })

    it("handles invalid date", () => {
        const result = time.localizeTime("sdfsdfsdfsdf", false, now)
        expect(result).toBeTruthy()
        expect(result.date).toBe("Invalid Date")
        expect(result.time).toBe("Invalid Date")
    })
})

describe("timeAgoString", () => {
    it("handles invalid date", () => {
        const result = time.timeAgoString("sdfsdfsdfsdf")
        expect(result).toBe("")
    })

    it("handles string with valid date", () => {
        const result = time.timeAgoString(yesterdayStr)
        expect(result).not.toBe("")
        expect(result).toBeTruthy()
    })
})
